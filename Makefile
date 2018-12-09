# .PHONY: serve
# serve:
# 	docker run -it -v C:\Users\aaron\source\Jemurai\web_application_threat_intelligence:/repo -p 9000:9000 store/gitpitch/desktop:pro

.PHONY: serve
serve:
	docker run -it -v /home/abedra/jemurai/web_application_threat_intelligence:/repo -p 9000:9000 store/gitpitch/desktop:pro
