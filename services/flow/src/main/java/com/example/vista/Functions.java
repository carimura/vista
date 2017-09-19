package com.example.vista;

import com.example.vista.messages.*;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fnproject.fn.api.Headers;
import com.fnproject.fn.api.flow.FlowFuture;
import com.fnproject.fn.api.flow.Flows;
import com.fnproject.fn.api.flow.HttpMethod;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;

/**
 * Created on 12/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
public class Functions {
    private static final ObjectMapper objectMapper = new ObjectMapper();
    private static final Logger log = LoggerFactory.getLogger(Functions.class);


    public static FlowFuture<ScrapeResp> runScraper(ScrapeReq req) {
        return wrapJsonFunction("./scraper", req, ScrapeResp.class);
    }

    public static FlowFuture<DetectPlateResp> detectPlates(DetectPlateReq req) {
        return wrapJsonFunction("./detect-plates", req, DetectPlateResp.class);
    }

    public static FlowFuture<DrawResp> drawRectangles(DrawReq req) {
        return wrapJsonFunction("./draw",req, DrawResp.class);
    }

    //
//    public static FlowFuture<DetectPlateResp> detectFaces(DetectPlateReq req) {
//        return wrapJsonFunction("./detect-faces",req,DetectPlateResp.class);
//    }


    public static FlowFuture<Void> postMessageToSlack(String channel,String message) {
        SlackRequest req = new SlackRequest();
        SlackMessage upload = new SlackMessage(message);
        req.message = upload;
        req.channel = channel;
        return wrapJsonFunction("./post-slack", req);

    }


    public static FlowFuture<Void> postImageToSlack(String channel,String url, String filename, String title, String initial_comment) {
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

    private static <RespT> FlowFuture<RespT> wrapJsonFunction(String name, Object input, Class<RespT> result) {
        log.info("Calling {} with {}", name, input);
        return Flows.currentFlow().invokeFunction(name, HttpMethod.POST, Headers.emptyHeaders(), toJson(input))
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
        log.info("Calling {} with {}", name, input);
        return Flows.currentFlow().invokeFunction(name, HttpMethod.POST, Headers.emptyHeaders(), toJson(input))
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
