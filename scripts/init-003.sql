\connect informatics

INSERT INTO coshh.chemical VALUES
                               (1, '123-45-6', 'Chemical one', '1', 'liquid', '100', CURRENT_DATE, CURRENT_DATE + 50, 'https://www.google.com', 'https://www.google.com', 'Lab 1', 'Cupboard 1', 'Owner 1', 'Shelf', 'false'),
                               (2, NULL, 'Chemical two', '2', 'liquid', '100', CURRENT_DATE, CURRENT_DATE + 50, 'https://www.google.com', 'https://www.google.com', 'Lab 1', 'Cupboard 2', 'Owner 2', '+4', 'false'),
                               (3, '345-67-8', 'Chemical three', NULL, 'solid', '100', CURRENT_DATE, CURRENT_DATE + 1, 'https://www.google.com', 'https://www.google.com', 'Lab 2', 'Cupboard 1', 'Owner 1', 'Shelf', 'false'),
                               (4, '456-78-9', 'Chemical four', '4', NULL, '100', CURRENT_DATE, CURRENT_DATE + 50, 'https://www.google.com', 'https://www.google.com', 'Lab 2', 'Cupboard 1a', 'Owner 3', '-20', 'false'),
                               (5, '567-89-0', 'Chemical five', '5', 'solid', NULL, CURRENT_DATE, CURRENT_DATE, 'https://www.google.com', 'https://www.google.com', 'Lab 3', 'Cupboard 1', 'Owner 3', 'Shelf', 'false'),
                               (6, '678-90-1', 'Chemical six', '6', 'liquid', '100', NULL, CURRENT_DATE + 50, 'https://www.google.com', 'https://www.google.com', 'Lab 3', 'Cupboard 1b', 'Owner 4', '+4', 'false'),
                               (7, '789-01-2', 'Chemical seven', '7', 'liquid', '100', CURRENT_DATE, NULL, 'https://www.google.com', 'https://www.google.com', 'Lab 4', 'Cupboard 1', 'Owner 4', '-80', 'false'),
                               (8, '890-12-3', 'Chemical eight', '8', 'solid', '100', CURRENT_DATE, CURRENT_DATE + 50, NULL, 'https://www.google.com', 'Lab 4', 'Cupboard 2', 'Owner 4', 'Shelf', 'false'),
                               (9, '890-12-4', 'Chemical nine', '9', 'liquid', '100', CURRENT_DATE, CURRENT_DATE + 50, 'https://www.google.com', NULL, 'Lab 5', 'Cupboard 3', 'Owner 5', '+4', 'false'),
                               (10, '890-12-5', 'Chemical ten', '10', 'liquid', '100', CURRENT_DATE, CURRENT_DATE + 50, 'https://www.google.com', 'https://www.google.com', NULL, 'Cupboard 4', 'Owner 1', 'Shelf', 'false'),
                               (11, '901-23-4', 'Chemical eleven', '11', 'solid', '100', CURRENT_DATE, CURRENT_DATE - 50, 'https://www.google.com', 'https://www.google.com', 'Lab 6', NULL, 'Owner 2', 'Shelf', 'false'),
                               (12, '012-34-5', 'Chemical twelve', '12', 'solid', '100', CURRENT_DATE, CURRENT_DATE + 50, 'https://www.google.com', 'https://www.google.com', 'Lab 6', 'Cupboard 1', NULL, '+4', 'false'),
                               (13, '111-11-1', 'Chemical thirteen', '13', 'liquid', '100', CURRENT_DATE, CURRENT_DATE + 50, 'https://www.google.com', 'https://www.google.com', 'Lab 7', 'Cupboard 12', 'Owner 4', '-20', 'false'),
                               (14, '222-22-2', 'Chemical fourteen', '14', 'solid', '100', CURRENT_DATE, CURRENT_DATE + 50, 'https://www.google.com', 'https://www.google.com', 'Lab 7', 'Cupboard 7', 'Owner 2', '-80', 'true');

INSERT INTO coshh.chemical_to_hazard VALUES
                                         (1, 'Explosive'),
                                         (2, 'Flammable'),
                                         (3, 'Oxidising'),
                                         (4, 'Corrosive'),
                                         (5, 'Acute toxicity'),
                                         (6, 'Hazardous to the environment'),
                                         (7, 'Health hazard/Hazardous to the ozone layer'),
                                         (8, 'Serious health hazard'),
                                         (9, 'Gas under pressure'),
                                         (10, 'Unknown'),
                                         (12, 'Explosive'),
                                         (12, 'Flammable'),
                                         (13, 'Oxidising'),
                                         (13, 'Corrosive'),
                                         (13, 'Acute toxicity'),
                                         (14, 'Hazardous to the environment'),
                                         (14, 'Health hazard/Hazardous to the ozone layer'),
                                         (14, 'Serious health hazard'),
                                         (14, 'Gas under pressure')