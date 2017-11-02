package com.example.vista.messages;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

/**
 * Created on 12/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public class ScrapeReq {
    public String query;
    public int num = 5;
    public int page;
}
