package com.example.vista.messages;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

/**
 * Created on 12/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
public class ScrapeResp implements Serializable {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public static class ScrapeResult implements Serializable {
        public String id;
        public String image_url;
    }

    public List<ScrapeResult> result = new ArrayList<>();
}
