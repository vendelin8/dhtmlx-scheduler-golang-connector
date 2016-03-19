window.onload = function() {
    scheduler.config.xml_date='%Y-%m-%d %H:%i';
    scheduler.config.first_hour = 9;
    scheduler.config.last_hour = 22;
    scheduler.config.multi_day = true;
    scheduler.config.date_step = '5';
    //scheduler.config.readonly = true;
    scheduler.templates.event_class = function(s, e, ev) {
        return ev.custom ? 'custom' : '';
    };
    scheduler.init('cntScheduler', new Date(), 'week');
    scheduler.setLoadMode('week');
    scheduler.load('/connector', 'json');
    var dp = new dataProcessor('/connector');
    dp.init(scheduler);
};
