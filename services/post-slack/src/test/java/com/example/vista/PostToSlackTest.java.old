
package com.example.vista;

import com.fnproject.fn.testing.FnTestingRule;
import org.junit.*;


@Ignore
public class PostToSlackTest{

    @Rule
    public final  FnTestingRule fn = FnTestingRule.createDefault();


    @Before
    public void setup(){
        String token = System.getenv("TEST_SLACK_API_TOKEN");
        Assume.assumeTrue("No slack testing token specified in environment", token!=null);
        fn.setConfig("SLACK_API_TOKEN",token);
    }

    @Test
    public void testPostMessageToSlack(){

        fn.givenEvent().withBody("{\n" +
                "  \"channel\": \"demostream\",\n" +
                "  \"message\": {\n" +
                "    \"text\": \"hello\"\n" +
                "  }\n" +
                "}").enqueue();

        fn.thenRun(PostToSlack.class,"postToSlack");
        Assert.assertEquals(200, fn.getOnlyResult().getStatus());
    }


    @Test
    public void testPostImageToSlack(){
        fn.givenEvent().withBody("{\n" +
                "  \"channel\": \"demostream\",\n" +
                "  \"upload\": {\n" +
                "    \"url\": \"https://farm3.staticflickr.com/2175/5714544755_e5dc8e6ede_b.jpg\",\n" +
                "    \"type\": \"auto\",\n" +
                "    \"filename\": \"movie.gif\",\n" +
                "    \"title\": \"Submarine\",\n" +
                "    \"initial_comment\": \"look an exploding submarine\"\n" +
                "  }\n" +
                "}").enqueue();

        fn.thenRun(PostToSlack.class,"postToSlack");
        Assert.assertEquals(200, fn.getOnlyResult().getStatus());

    }
}