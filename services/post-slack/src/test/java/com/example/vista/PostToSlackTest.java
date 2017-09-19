
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
                "  \"channel\": \"general\",\n" +
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
                "  \"channel\": \"general\",\n" +
                "  \"upload\": {\n" +
                "    \"url\": \"http://www.zoot.org.uk/wp/wp-content/uploads/2015/01/dragknife2-150x150.jpeg\",\n" +
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