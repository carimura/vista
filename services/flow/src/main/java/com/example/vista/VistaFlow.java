package com.example.vista;

import com.example.vista.messages.DetectPlateReq;
import com.example.vista.messages.DrawReq;
import com.example.vista.messages.ScrapeReq;
import com.fnproject.fn.api.flow.FlowFuture;
import com.fnproject.fn.api.flow.Flows;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;

import static com.example.vista.Functions.*;

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
                    log.info("Got  {} images from the scraper ", resp.result.size());
                    List<FlowFuture<?>> pendingTasks = new ArrayList<>();

                    resp.result.forEach(scrapeResult -> {

                        log.info("starting detection on {}", scrapeResult.image_url);

                        String id = scrapeResult.id;

                        FlowFuture<?> processTask = detectPlates(new DetectPlateReq(scrapeResult.image_url, "us"))

                                .thenCompose((plateResp) -> {

                                    if (!plateResp.got_plate) {
                                        log.info("No plates in {}", scrapeResult.image_url);
                                        return Flows.currentFlow().completedValue(null);
                                    }
                                    log.info("Got plate {} in {}", plateResp.plate, scrapeResult.image_url);
                                    return Functions
                                            .drawRectangles(new DrawReq(id, scrapeResult.image_url, plateResp.rectangles))
                                            .thenCompose((drawResp) ->
                                                    postImageToSlack(slackChannel, drawResp.image_url, null, "Found Plate" + plateResp.plate, "Have you seen this car?"));

                                });

                        pendingTasks.add(processTask);
                    });

                    return Flows.currentFlow().allOf(pendingTasks.toArray(new FlowFuture[pendingTasks.size()]));
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
