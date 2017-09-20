package com.example.vista;

import com.example.vista.messages.DetectPlateReq;
import com.example.vista.messages.DrawReq;
import com.example.vista.messages.ScrapeReq;
import com.example.vista.messages.ScrapeResp;
import com.fnproject.fn.api.flow.FlowFuture;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.stream.Collectors;

import static com.example.vista.Functions.*;
import static com.fnproject.fn.api.flow.Flows.currentFlow;

/**
 * FnFlow function that composes : several functions together
 */
public class VistaFlow {
    static {
        System.setProperty("org.slf4j.simpleLogger.defaultLogLevel", "debug");
    }

    private static final Logger log = LoggerFactory.getLogger(VistaFlow.class);


    public void handleRequest(ScrapeReq input) throws Exception {

        log.info("Got request {} {}", input.query, input.num);

        String slackChannel = "general";

        postMessageToSlack(slackChannel, String.format("About to start scraping for images containing \"%s\"", input.query)).get();

        runScraper(input)
                .thenCompose(resp -> {

                    List<ScrapeResp.ScrapeResult> results = resp.result;

                    log.info("Got {} images from the scraper ", results.size());

                    List<FlowFuture<?>> pendingTasks = results
                            .stream()
                            .map(scrapeResult -> {
                                log.info("starting image detection on {}", scrapeResult.image_url);

                                String id = scrapeResult.id;

                                return detectPlates(new DetectPlateReq(scrapeResult.image_url, "us")).thenCompose((plateResp) -> {
                                    if (!plateResp.got_plate) {
                                        log.info("No plates found in {}", scrapeResult.image_url);
                                        return currentFlow().completedValue(null);
                                    }
                                    log.info("Got plate {} in {}", plateResp.plate, scrapeResult.image_url);
                                    return Functions
                                            .drawRectangles(new DrawReq(id, scrapeResult.image_url, plateResp.rectangles))
                                            .thenCompose((drawResp) ->
                                                    currentFlow().allOf(
                                                            postAlertToTwitter(drawResp.image_url, plateResp.plate),
                                                            postImageToSlack(slackChannel,
                                                                    drawResp.image_url,
                                                                    null,
                                                                    "Found Plate" + plateResp.plate,
                                                                    "Have you seen this car?")));

                                });

                            }).collect(Collectors.toList());

                    return currentFlow().allOf(pendingTasks.toArray(new FlowFuture[pendingTasks.size()]));
                })
                .whenComplete((v, throwable) -> {
                    if (throwable != null) {
                        log.info("Scraping completed with at least one error", throwable);
                        postMessageToSlack(slackChannel, "Something went wrong!" + throwable.getMessage());

                    } else {
                        log.info("Scraping completed successfully");

                        postMessageToSlack(slackChannel, "Finished scraping");

                    }
                });

    }

}
