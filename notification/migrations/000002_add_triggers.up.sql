CREATE OR REPLACE FUNCTION check_subscription_requirements()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.type = 'telegram' THEN
        IF (SELECT telegram_id FROM users WHERE user_id = NEW.user_id) IS NULL THEN
            RAISE EXCEPTION 'User must have a telegram_id to subscribe to telegram notifications';
        END IF;
    ELSIF NEW.type = 'email' THEN
        IF (SELECT email FROM users WHERE user_id = NEW.user_id) IS NULL THEN
            RAISE EXCEPTION 'User must have an email to subscribe to email notifications';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_subscription
BEFORE INSERT ON subscriptions
FOR EACH ROW
EXECUTE FUNCTION check_subscription_requirements();
