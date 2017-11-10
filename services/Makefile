APP?=vista

deploy:
	cd draw; fn deploy --app ${APP}
	cd publish; fn deploy --app ${APP}
	cd detect-plates; fn deploy --app ${APP}
	cd alert; fn deploy --app ${APP}
	cd scraper; fn deploy --app ${APP}
	cd flow; fn deploy --app ${APP}
	cd post-slack; fn deploy --app ${APP}

deploy-local:
	cd draw; fn deploy --local --app ${APP}
	cd publish; fn deploy --local --app ${APP}
	cd detect-plates; fn deploy --local --app ${APP}
	cd alert; fn deploy --local --app ${APP}
	cd scraper; fn deploy --local --app ${APP}
	cd flow; fn deploy --local --app ${APP}
	cd post-slack; fn deploy --local --app ${APP}
