# FnFlow vista function 

This creates an asynchronous flow across different parts of the vista app. 

Rather than chaining functions together within functions the flow call creates a FnFlow that does the following: 

* posts a message on slack 
* Calls the scraper to get a list of images
* for each image it: 
   * runs detect 
   * runs draw 
   * then in parallel it: 
       * Posts to /alert
       * Posts to /slack 
       * Joins on the completion of those
* Joins on all images being processes
* Posts to slack 



 