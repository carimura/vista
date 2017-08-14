init: 
	cd draw; fn init --type async $(dockerid)/draw;
	cd publish; fn init --type async $(dockerid)/publish;
	cd detect-plates; fn init --type async $(dockerid)/detect-plates;
	cd detect-faces; fn init --type async $(dockerid)/detect-faces;
	cd alert; fn init --type async $(dockerid)/alert;
	cd scraper; fn init $(dockerid)/scraper;


deploy:
	cd draw; fn deploy myapp
	cd publish; fn deploy myapp
	cd detect-plates; fn deploy myapp
	cd detect-faces; fn deploy myapp
	cd alert; fn deploy myapp
	cd scraper; fn deploy myapp
