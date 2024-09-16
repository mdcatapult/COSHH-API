CREATE OR REPLACE FUNCTION export_chemical_data(file_path TEXT)
RETURNS VOID AS $$
BEGIN
    COPY (
        SELECT
            coshh.chemical.*,
            coshh.chemical_to_hazard.*
        FROM
            coshh.chemical
        JOIN
            coshh.chemical_to_hazard ON coshh.chemical.id = coshh.chemical_to_hazard.id
    ) TO file_path WITH CSV HEADER;
END; 
$$ LANGUAGE plpgsql;

SELECT export_chemical_data('/path/to/download/chemical/to/azure/data.csv');

