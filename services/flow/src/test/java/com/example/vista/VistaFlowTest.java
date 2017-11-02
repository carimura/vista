package com.example.vista;

import com.fnproject.fn.testing.*;
import org.junit.*;

public class VistaFlowTest {

    @Rule
    public final FnTestingRule testing = FnTestingRule.createDefault();

    @Test
    public void shouldFindImages() {

        testing.givenFn("./scraper").withResult(("{\"result\":[{\"id\":\"34053257076\",\"image_url\":\"https://farm3.staticflickr.com/2883/34053257076_2911069a6d_c.jpg\"}]}").getBytes());
        testing.givenFn("./detect-plates").withResult(("{\"got_plate\":true,\"rectangles\":[{\"startx\":834,\"starty\":702,\"endx\":1022,\"endy\":783}],\"plate\":\"D33M016\"}").getBytes());
        testing.givenFn("./draw").withAction((b)->{
            System.err.println("got drawRectangles");
            return "{\"image_url\":\"http://example.com/image.png\"}".getBytes();
        });
        testing.givenFn("./alert").withResult(("OK").getBytes());
        testing.givenFn("./post-slack").withResult(("").getBytes());


        testing.givenEvent().withBody("{\"query\": \"license plate car usa\"}")
                .enqueue();

        testing.thenRun(VistaFlow.class, "handleRequest");

        FnResult result = testing.getOnlyResult();
        Assert.assertEquals("", result.getBodyAsString());
    }

}