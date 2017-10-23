CREATE TABLE restaurants(
	id                        serial,
	establishmentID           int UNIQUE,	
	establishmentName         varchar,
	establishmentType         varchar,
	establishmentAddress      varchar,
	establishmentStatus       varchar,
	minimumInspectionsPerYear int
);

CREATE TABLE inspections(
	establishmentID           int,
	inspectionID              int UNIQUE,
	infractionDetails         varchar,
	inspectionDate            varchar,
	severity                  varchar,
	action                    varchar,
	courtOutcome              varchar,
	amountFinded              varchar
);
