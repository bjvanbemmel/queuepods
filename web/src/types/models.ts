export enum Events {
  QUEUE_EMPTY           = "queue_empty",
  QUEUE_FULL            = "queue_full",
  QUEUE_ALMOST_EMPTY    = "queue_almost_empty",
  QUEUE_ALMOST_FULL     = "queue_almost_full",
  POPULATION_MONITORING = "population_monitoring",
};

export enum Params {
  EVENTS =      "events",
  LIMIT =       "limit",
  ATTRACTIONS = "attractions",
  FROM =        "from",
};

export interface Param {
  name:  Params,
  value: string,
};

export interface Population {
  attraction: string,
  population: number,
  state:      string,
  capacity:   number,
};

export interface Message {
  attraction: string,
  timestamp:  Date,
  event:      Events,
  value:      string,
};
