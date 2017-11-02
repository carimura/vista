package com.example.vista.messages;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import java.io.Serializable;

/**
 * Created on 12/09/2017.
 * <p>
 * (c) 2017 Oracle Corporation
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public class Rect implements Serializable {
    public int startx;
    public int starty;
    public int endx;
    public int endy;
}
