package com.example.vista;

import com.fnproject.fn.testing.*;
import org.junit.*;
import com.fnproject.fn.testing.FnTestingRule;
import com.fnproject.fn.testing.flow.FlowTesting;

public class VistaFlowTest {

    @Rule
    public final FnTestingRule fn = FnTestingRule.createDefault();
    private final FlowTesting flow = FlowTesting.create(fn);

    @Test
    public void shouldFindImages() {

        flow.givenFn("./scraper").withResult(("{\"result\":[{\"id\":\"34053257076\",\"image_url\":\"https://farm3.staticflickr.com/2883/34053257076_2911069a6d_c.jpg\"}]}").getBytes());
        flow.givenFn("./detect-plates").withResult(("{\"got_plate\":true,\"rectangles\":[{\"startx\":834,\"starty\":702,\"endx\":1022,\"endy\":783}],\"plate\":\"D33M016\"}").getBytes());
        flow.givenFn("./draw").withAction((b)->{
            System.err.println("got drawRectangles");
            return "{\"image_url\":\"http://example.com/image.png\"}".getBytes();
        });
        flow.givenFn("./alert").withResult(("OK").getBytes());
        flow.givenFn("./post-slack").withResult(("").getBytes());


        fn.givenEvent().withBody("{\"query\": \"license plate car usa\"}")
                .enqueue();

        fn.thenRun(VistaFlow.class, "handleRequest");

        FnResult result = fn.getOnlyResult();
        Assert.assertEquals("", result.getBodyAsString());
    }

}