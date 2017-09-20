package com.example.vista.messages;

import java.util.List;

/**
 * Created on 12/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
public class DrawReq {
    public final String id;
    public final String image_url;
    public final List<Rect> rectangles;
    public final String resize;

    public DrawReq(String id, String image_url, List<Rect> rectangles, String resize) {
        this.id = id;
        this.image_url = image_url;
        this.rectangles = rectangles;
        this.resize = resize;
    }
}
