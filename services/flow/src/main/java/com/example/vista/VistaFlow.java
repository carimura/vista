package com.example.vista;

import com.example.vista.messages.*;
import com.fnproject.fn.api.FnConfiguration;
import com.fnproject.fn.api.RuntimeContext;
import com.fnproject.fn.api.flow.FlowFuture;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.Serializable;
import java.util.List;
import java.util.stream.Collectors;

import static com.example.vista.Slack.postMessageToSlack;
import static com.fnproject.fn.api.flow.Flows.currentFlow;
import com.fnproject.fn.runtime.flow.FlowFeature;
import com.fnproject.fn.api.FnFeature;

/**
 * FnFlow Vista function
 */
@FnFeature(FlowFeature.class)
public class VistaFlow implements Serializable {
  static {
    System.setProperty("org.slf4j.simpleLogger.defaultLogLevel", "debug");
  }

  public void handleRequest(ScrapeReq input) throws Exception {

    log.info("Got request {} {}", input.query, input.num);
    postMessageToSlack(slackFuncID, slackChannel, String.format("About to start scraping for images containing \"%s\"", input.query)).get();
    FlowFuture<ScrapeResp> scrapes = currentFlow().invokeFunction(scraperFuncID, input, ScrapeResp.class);

    scrapes.thenCompose(resp -> {

      List<ScrapeResp.ScrapeResult> results = resp.result;
      log.info("Got {} images from the scraper ", results.size());
      List<FlowFuture<?>> pendingTasks = results
          .stream()
          .map(scrapeResult -> {
            log.info("starting image detection on {}", scrapeResult.image_url);

            String id = scrapeResult.id;
            return currentFlow()
                .invokeFunction(detectPlatesFuncID, new DetectPlateReq(scrapeResult.image_url, "us"), DetectPlateResp.class)
                .thenCompose((plateResp) -> {

                  if (!plateResp.got_plate) {
                    log.info("No plates found in {}", scrapeResult.image_url);
                    // bug
                    return currentFlow().completedValue(null);
                  }

                  log.info("Got plate {} in {}", plateResp.plate, scrapeResult.image_url);
                  return currentFlow()
                      .invokeFunction(drawFuncID, new DrawReq(id, scrapeResult.image_url, plateResp.rectangles,"300x300"), DrawResp.class)
                      .thenCompose((drawResp) -> {
                            // Finally when the image is rendered  post an alert to twitter and slack in parallel
                            log.info("Got draw response {} ", drawResp.image_url);
                            return currentFlow().allOf(
                                currentFlow().invokeFunction(alertFuncID, new AlertReq(plateResp.plate, drawResp.image_url)),
                                Slack.postImageToSlack(slackFuncID, slackChannel,
                                    drawResp.image_url,
                                    "plate",
                                    "Found plate: " + plateResp.plate,
                                    "Have you seen this car?"));
                          }
                      );

                });
          }).collect(Collectors.toList());

      return currentFlow()
          .allOf(pendingTasks.toArray(new FlowFuture[pendingTasks.size()]));

    }).whenComplete((v, throwable) -> {
      if (throwable != null) {
        log.info("Scraping completed with at least one error", throwable);
        postMessageToSlack(slackFuncID, slackChannel, "Something went wrong!" + throwable.getMessage());

      } else {
        log.info("Scraping completed successfully");
        postMessageToSlack(slackFuncID, slackChannel, "Finished scraping");
      }
    });
  }


  private static final Logger log = LoggerFactory.getLogger(VistaFlow.class);
  private static String slackChannel = "demostream";
  // func IDs are necessary
  private static String slackFuncID = null;
  private static String alertFuncID = null;
  private static String scraperFuncID = null;
  private static String detectPlatesFuncID = null;
  private static String drawFuncID = null;

  @FnConfiguration
  private static void configure(RuntimeContext ctx) {
    ctx.getConfigurationByKey("SLACK_CHANNEL").ifPresent((c) -> slackChannel = c);
    // func IDs are really necessary
    ctx.getConfigurationByKey("POST_SLACK_FUNC_ID").ifPresent((c) -> slackFuncID = c);
    ctx.getConfigurationByKey("SCRAPER_FUNC_ID").ifPresent((c) -> scraperFuncID = c);
    ctx.getConfigurationByKey("DETECT_PLATES_FUNC_ID").ifPresent((c) -> detectPlatesFuncID = c);
    ctx.getConfigurationByKey("ALERT_FUNC_ID").ifPresent((c) -> alertFuncID = c);
    ctx.getConfigurationByKey("DRAW_FUNC_ID").ifPresent((c) -> drawFuncID = c);
  }
}
