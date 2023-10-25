-- Description: This script creates the audit trigger function and the audit log table
CREATE TABLE coshh.audit_coshh_logs(
    audit_coshh_loguser_id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username text NOT NULL, -- current_user in trigger function
    table_name text NOT NULL, -- TG_TABLE_NAME in trigger function
    operation text NOT NULL, -- TG_OP in trigger function
    last_updated TIMESTAMP NOT NULL DEFAULT NOW(), -- time of event
    row_data jsonb NOT NULL -- coshh database records converted to jsonb
)

-- This trigger function is used to audit all changes to the database
CREATE OR REPLACE FUNCTION coshh_audit_triggerfunction()  
    RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP IN ('UPDATE', 'DELETE', 'INSERT', 'TRUNCATE') THEN
        INSERT INTO coshh.audit_coshh_logs (
            username, 
            table_name, 
            operation ,
            last_updated,
            row_data
        )
        VALUES (
            session_user, 
            TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME, 
            TG_OP,
            NOW(),
            jsonb_build_object(
                'old', row_to_json(OLD.*),
                'new', row_to_json(NEW.*)
            )
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- The trigger to call the audit trigger function
CREATE TRIGGER audit_coshh_logs
    AFTER INSERT OR UPDATE OR DELETE ON coshh.chemical
    FOR EACH ROW 
EXECUTE PROCEDURE coshh_audit_triggerfunction();


-- disallow truncate
CREATE TRIGGER audit_coshh_logs_truncate
    BEFORE TRUNCATE ON coshh.chemical
    FOR EACH STATEMENT
EXECUTE PROCEDURE coshh_audit_triggerfunction();


-- A view to show the audit log with just id and last_updated
CREATE OR REPLACE VIEW coshh.audit_coshh_log_views AS
    SELECT coshh_dev.id,
           max(audit_coshh.last_updated) AS last_updated
    FROM coshh.chemical coshh_dev
    LEFT JOIN coshh.audit_coshh_logs audit_coshh
        ON audit_coshh.audit_coshh_loguser_id = coshh_dev.id
    GROUP BY coshh_dev.id;


-- A view function to make it updatedable so any type of operation performed on the view displays the last modified time of the record
CREATE OR REPLACE FUNCTION audit_coshh_log_views_function()  
RETURNS TABLE (audit_coshh_viewer_id int, date_updated TIMESTAMP , columns_updated jsonb) AS $$
DECLARE
    row_record RECORD;
    new_data jsonb;
    old_data jsonb;
	new_key text;
	new_value text;
	old_key text; 
	old_value text;
BEGIN
    -- loop through the audit table and compare the new and old data
    FOR row_record IN SELECT audit_coshh_loguser_id, last_updated, row_data FROM coshh.audit_coshh_logs 
    WHERE coshh.audit_coshh_logs.operation = 'UPDATE'
    LOOP
        columns_updated := '{}'::jsonb;

        -- extract the 'new' and 'old' object key words from the jsonb object
        new_data :=  row_record.row_data->'new';
        old_data :=  row_record.row_data->'old';

        -- loop through the new data and compare it to the old data
        FOR new_key, new_value IN SELECT * FROM jsonb_each_text(new_data)
        LOOP
			FOR old_key, old_value IN SELECT * FROM jsonb_each_text(old_data) 
			LOOP
				IF new_key = old_key OR new_value = 'id'
                THEN
                    IF new_value != old_value
                    THEN
                        columns_updated = jsonb_set(columns_updated, ARRAY[new_key], to_jsonb(new_value));
                    END IF;
                END IF;
			END LOOP;
        END LOOP;
        IF columns_updated != '{}'::jsonb
            THEN
                audit_coshh_viewer_id := row_record.audit_coshh_loguser_id;
                date_updated := row_record.last_updated;
            RETURN NEXT;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;


-- The trigger to call the audit trigger function
CREATE TRIGGER audit_trigger_coshh_logs
    AFTER INSERT OR UPDATE OR DELETE ON coshh.audit_coshh_logs
    FOR EACH ROW 
EXECUTE PROCEDURE audit_coshh_log_views_function();


-- a view to access the function's output
CREATE OR REPLACE VIEW coshh.audit_coshh_log_views AS
    SELECT * FROM audit_coshh_log_views_function();


-- A generic cancel trigger function to prevent changes to the database
CREATE OR REPLACE FUNCTION cancel_triggerfunction()  
    RETURNS TRIGGER AS $$
BEGIN
    IF TG_WHEN = 'AFTER' THEN
        RAISE EXCEPTION 'You are not allowed to % ROWS IN %.%',       
                              TG_OP, 
                              TG_TABLE_SCHEMA,
                              TG_TABLE_NAME;
    END IF;
    RAISE NOTICE '% ON ROWS IN %.% IS WONT BE EXECUTED', 
                              TG_OP, 
                              TG_TABLE_SCHEMA,
                              TG_TABLE_NAME;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;


-- disallow insert, update, delete
CREATE TRIGGER cancel_trigger
    BEFORE INSERT OR UPDATE OR DELETE ON coshh.chemical
    FOR EACH ROW 
EXECUTE PROCEDURE cancel_triggerfunction();


-- disallow truncate
CREATE TRIGGER cancel_trigger
    BEFORE TRUNCATE ON coshh.chemical
    FOR EACH STATEMENT
EXECUTE PROCEDURE cancel_triggerfunction();