package com.example.vista;

import com.example.vista.messages.SlackMessage;
import com.example.vista.messages.SlackRequest;
import com.example.vista.messages.SlackUpload;
import com.fnproject.fn.api.flow.FlowFuture;
import com.fnproject.fn.api.flow.Flows;
import com.fnproject.fn.api.flow.HttpResponse;

/**
 * Slack stubs - these are methods that create typed, asynchronous calls to slack functions
 * <p>
 * Each method returns a FlowFuture that yeilds the result of the function call (or an error if one occured)
 * <p>
 * Created on 12/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
public class Slack {


    public static FlowFuture<HttpResponse> postMessageToSlack(String slackFuncID, String channel, String message) {
        SlackRequest req = new SlackRequest();
        req.message = new SlackMessage(message);
        req.channel = channel;
        return Flows.currentFlow().invokeFunction(slackFuncID, req);
    }


    public static FlowFuture<HttpResponse> postImageToSlack(String slackFuncID, String channel, String url, String filename, String title, String initial_comment) {
        SlackRequest req = new SlackRequest();
        SlackUpload upload = new SlackUpload();
        upload.filename = filename;
        upload.title = title;
        upload.initial_comment = initial_comment;
        upload.url = url;
        req.upload = upload;
        req.channel = channel;
        return Flows.currentFlow().invokeFunction(slackFuncID, req);
    }
}
