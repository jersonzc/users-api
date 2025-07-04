CREATE FUNCTION update_updated_at_column() RETURNS trigger AS $$
    BEGIN
        NEW.updated_at = CURRENT_TIMESTAMP;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at_trigger
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
