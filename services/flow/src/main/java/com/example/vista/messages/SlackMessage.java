package com.example.vista.messages;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

/**
 * Created on 19/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public class SlackMessage {
    public final String text;

    public SlackMessage(String text) {
        this.text = text;
    }
}
