package com.example.vista.messages;

import java.io.Serializable;
import java.util.List;

/**
 * Created on 12/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
public class DetectPlateResp implements Serializable {
    public boolean got_plate;
    public String plate;

    public List<Rect> rectangles;

}
