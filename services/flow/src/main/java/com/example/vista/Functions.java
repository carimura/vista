package com.example.vista;

import com.example.vista.messages.*;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fnproject.fn.api.Headers;
import com.fnproject.fn.api.flow.FlowFuture;
import com.fnproject.fn.api.flow.Flows;
import com.fnproject.fn.api.flow.HttpMethod;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

/**
 *
 * Function stubs - these are methods that create typed, asynchronous calls to the Faas
 *
 * Each method returns a FlowFuture that yeilds the result of the function call (or an error if one occured)
 *
 * Created on 12/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
public class Functions {
    private static final ObjectMapper objectMapper = new ObjectMapper().configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
    private static final Logger log = LoggerFactory.getLogger(Functions.class);


    public static FlowFuture<ScrapeResp> runScraper(ScrapeReq req) {
        return wrapJsonFunction("./scraper", req, ScrapeResp.class);
    }

    public static FlowFuture<DetectPlateResp> detectPlates(DetectPlateReq req) {
        return wrapJsonFunction("./detect-plates", req, DetectPlateResp.class);
    }

    public static FlowFuture<DrawResp> drawRectangles(DrawReq req) {
        return wrapJsonFunction("./draw", req, DrawResp.class);
    }


    public static FlowFuture<Void> postMessageToSlack(String channel, String message) {
        SlackRequest req = new SlackRequest();
        req.message = new SlackMessage(message);
        req.channel = channel;
        return wrapJsonFunction("./post-slack", req);

    }


    public static FlowFuture<Void> postImageToSlack(String channel, String url, String filename, String title, String initial_comment) {
        SlackRequest req = new SlackRequest();
        SlackUpload upload = new SlackUpload();
        upload.filename = filename;
        upload.title = title;
        upload.initial_comment = initial_comment;
        upload.url = url;
        req.upload = upload;
        req.channel = channel;
        return wrapJsonFunction("./post-slack", req);

    }

    public static FlowFuture<Void> postAlertToTwitter(String url, String plate) {
        return wrapJsonFunction("./alert", new AlertReq(plate, url));
    }

    private static <RespT> FlowFuture<RespT> wrapJsonFunction(String name, Object input, Class<RespT> result) {
        byte[] bytes = toJson(input);

        log.info("Calling {} with {}:{}", name, input, new String(bytes));
        Map<String, String> headerMap = new HashMap<>();
        headerMap.put("NO_CHAIN", "true");

        return Flows.currentFlow().invokeFunction(name, HttpMethod.POST, Headers.fromMap(headerMap), bytes)
                .thenApply((httpResp) -> fromJson(httpResp.getBodyAsBytes(), result))
                .whenComplete((v, e) -> {
                    if (e != null) {
                        log.error("Got error from {} ", name, e);
                    } else {
                        log.info("Got response from {}: {}", name, v);
                    }
                });

    }


    private static FlowFuture<Void> wrapJsonFunction(String name, Object input) {
        byte[] bytes = toJson(input);
        log.info("Calling {} with {} : {}", name, input, new String(bytes));

        Map<String, String> headerMap = new HashMap<>();
        headerMap.put("NO_CHAIN", "true");

        return Flows.currentFlow().invokeFunction(name, HttpMethod.POST, Headers.fromMap(headerMap), bytes)
                .handle((v, e) -> {
                    if (e != null) {
                        log.error("Got error from {} ", name, e);
                    } else {
                        log.info("Got response from {}: {}", name, v);
                    }
                    return null;
                });

    }


    private static <T> T fromJson(byte[] data, Class<T> type) {
        try {
            return objectMapper.readValue(data, type);
        } catch (IOException e) {
            log.error("Failed to extract value to {} ", type, e);
            throw new RuntimeException(e);
        }
    }

    private static <T> byte[] toJson(T val) {
        try {
            return objectMapper.writeValueAsString(val).getBytes();
        } catch (IOException e) {
            log.error("Failed to wite {} to json ", val, e);
            throw new RuntimeException(e);
        }
    }


}
