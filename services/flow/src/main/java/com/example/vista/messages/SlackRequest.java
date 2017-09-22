package com.example.vista.messages;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

/**
 * Created on 19/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public class SlackRequest {
    public String channel;
    public SlackUpload upload;
    public SlackMessage message;
}
