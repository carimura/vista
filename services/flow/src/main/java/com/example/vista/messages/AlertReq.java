package com.example.vista.messages;

/**
 * Created on 12/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
public class AlertReq {
    public final String plate;
    public final String image_url;

    public AlertReq(String plate, String image_url) {
        this.plate = plate;
        this.image_url = image_url;
    }
}
