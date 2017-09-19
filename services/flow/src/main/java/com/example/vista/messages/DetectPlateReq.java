package com.example.vista.messages;

/**
 * Created on 12/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
public class DetectPlateReq {
    public final String image_url;
    public final String countrycode;

    public DetectPlateReq(String image_url, String countrycode) {
        this.image_url = image_url;
        this.countrycode = countrycode;
    }
}
