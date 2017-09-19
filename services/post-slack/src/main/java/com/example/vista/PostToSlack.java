package com.example.vista;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fnproject.fn.api.FnConfiguration;
import com.fnproject.fn.api.RuntimeContext;
import okhttp3.*;
import org.apache.commons.io.IOUtils;

import java.io.IOException;
import java.net.URI;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.TimeUnit;

public class PostToSlack {

    private static String slackToken;
    private static OkHttpClient client = new OkHttpClient.Builder().readTimeout(60, TimeUnit.SECONDS).build();
    private static ObjectMapper objectMapper = new ObjectMapper();

    @FnConfiguration
    public static void configure(RuntimeContext ctx) {
        slackToken = ctx.getConfigurationByKey("SLACK_API_TOKEN").orElseThrow(() -> new RuntimeException("Missing SLACK_API_TOKEN config variable"));
    }

    public static class Upload {
        public String url;
        public String type = "auto";
        public String filename = "";
        public String title = "";
        public String initial_comment = "";
    }

    public static class Message {
        public String text;
    }

    public static class SlackRequest {
        public String channel;
        public Upload upload;
        public Message message;
    }

    @JsonIgnoreProperties(ignoreUnknown = true)
    public static class SlackResponse {
        public boolean ok;
        public String error;
    }

    public void postToSlack(SlackRequest input) throws Exception {

        if (input.message != null) {

            sendRequest(new Request.Builder()
                    .url(new HttpUrl.Builder()
                            .scheme("https")
                            .host("slack.com")
                            .addPathSegments("api/chat.postMessage")
                            .addEncodedQueryParameter("text", URLEncoder.encode(input.message.text, StandardCharsets.UTF_8.toString()).replace("?", "%3F"))
                            .addQueryParameter("token", slackToken)
                            .addQueryParameter("channel", input.channel)
                            .addQueryParameter("as_user", "true")
                            .build())
                    .build());
        } else if (input.upload != null) {
            byte[] data = IOUtils.toByteArray(URI.create(input.upload.url));

            RequestBody requestBody = new MultipartBody.Builder()
                    .setType(MultipartBody.FORM)
                    .addFormDataPart("token", slackToken)
                    .addFormDataPart("channels", input.channel)
                    .addFormDataPart("filetype", input.upload.type)
                    .addFormDataPart("initial_comment", input.upload.initial_comment)
                    .addFormDataPart("title", input.upload.title)
                    .addFormDataPart("file", input.upload.filename, RequestBody.create(null, data))
                    .build();

            Request r = new Request.Builder()
                    .url(new HttpUrl.Builder()
                            .scheme("https")
                            .host("slack.com")
                            .addPathSegments("api/files.upload").build())
                    .post(requestBody)
                    .build();

            sendRequest(r);
        }
    }

    private void sendRequest(Request request) throws IOException {
        Response res = client.newCall(request).execute();
        System.err.println("Got response " + res);
        if (!res.isSuccessful()) {
            throw new RuntimeException("Invalid response : " + res.toString());
        }

        String result = res.body().string();
        System.err.println("Got result" + result);
        SlackResponse sm = objectMapper.readValue(result, SlackResponse.class);
        if (!sm.ok) {
            throw new RuntimeException("Error from slack API :" + sm.error);
        }
    }

}