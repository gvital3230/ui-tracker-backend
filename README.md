# Sample UI tracker application - Backend

This is the part of sample UI tracker service, based on Websocket communication

**While launching demo apps please be aware of Heroku [Free tier Dyno sleeping](https://devcenter.heroku.com/articles/free-dyno-hours#dyno-sleeping)
so it can be possibly be several seconds delay at initial load**

Service consists of three parts

### Public application

Small JS application which simulates user activity on some e-commerce page. User opens it, starts to scroll.
After some delay application tracks activity related to every product image and send it to backend server

* [Repository location](https://github.com/gvital3230/ui-tracker-public)
* [Demo app](https://ui-tracker-public.herokuapp.com) - to see, how it works, it is recommended to open it from multiply different browsers or devices

![](assets/howto.gif)

### Dashboard application

Shows current active sessions and images that users are seeing right now

* [Repository location](https://github.com/gvital3230/ui-tracker-dashboard)
* [Demo app](https://ui-tracker-dashboard.herokuapp.com)

### Backend application

Service that responsible for processing Websocket connection and build current active sessions map

* Repository location - (this one)
* Demo app - all settings are made in demo public/dashboard apps

