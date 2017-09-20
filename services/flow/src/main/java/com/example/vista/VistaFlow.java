package com.example.vista;

import com.example.vista.messages.DetectPlateReq;
import com.example.vista.messages.DrawReq;
import com.example.vista.messages.ScrapeReq;
import com.example.vista.messages.ScrapeResp;
import com.fnproject.fn.api.RuntimeContext;
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
    private static String slackChannel = "demostream";

    private static void configure(RuntimeContext ctx) {
        ctx.getConfigurationByKey("SLACK_CHANNEL").ifPresent((c) -> {
            slackChannel = c;
        });
    }

    public void handleRequest(ScrapeReq input) throws Exception {

        log.info("Got request {} {}", input.query, input.num);


        postMessageToSlack(slackChannel, String.format("About to start scraping for images containing \"%s\"", input.query)).get();

        // Run the scraper function to get images from Flickr
        runScraper(input)
                .thenCompose(resp -> {
                    // take the list of images  and start a detect job for each image in paralle
                    List<ScrapeResp.ScrapeResult> results = resp.result;

                    log.info("Got {} images from the scraper ", results.size());

                    List<FlowFuture<?>> pendingTasks = results
                            .stream()
                            .map(scrapeResult -> {
                                log.info("starting image detection on {}", scrapeResult.image_url);

                                String id = scrapeResult.id;

                                // detect plates on each image
                                return detectPlates(new DetectPlateReq(scrapeResult.image_url, "us")).thenCompose((plateResp) -> {
                                    if (!plateResp.got_plate) {
                                        log.info("No plates found in {}", scrapeResult.image_url);
                                        // bug
                                        return currentFlow().completedValue(null);
                                    }

                                    // if a plate was found trigger a job to draw the detected rectangles on the original image
                                    // Draw rectangles renders the image to an object store returns the URL of the image
                                    log.info("Got plate {} in {}", plateResp.plate, scrapeResult.image_url);
                                    return Functions
                                            .drawRectangles(new DrawReq(id, scrapeResult.image_url, plateResp.rectangles, "300x300"))
                                            .thenCompose((drawResp) -> {

                                                        // Finally when the image is rendered  post an alert to twitter and slack in parallel
                                                        log.info("Got draw response {} ", drawResp.image_url);
                                                        return currentFlow().allOf(
                                                                postAlertToTwitter(drawResp.image_url, plateResp.plate),
                                                                postImageToSlack(slackChannel,
                                                                        drawResp.image_url,
                                                                        "plate",
                                                                        "Found plate: " + plateResp.plate,
                                                                        "Have you seen this car?"));
                                                    }
                                            );

                                });

                                // Collect all of the tasks into a list
                            }).collect(Collectors.toList());

                    return currentFlow().allOf(pendingTasks.toArray(new FlowFuture[pendingTasks.size()]));
                })
                .whenComplete((v, throwable) -> {
                    // when all of the tasks are complete, join on them and send a final message to slack
                    if (throwable != null) {
                        // if an error occured we get a stack trace here.
                        log.info("Scraping completed with at least one error", throwable);
                        postMessageToSlack(slackChannel, "Something went wrong!" + throwable.getMessage());

                    } else {
                        log.info("Scraping completed successfully");
                        postMessageToSlack(slackChannel, "Finished scraping");

                    }
                });

    }

}
