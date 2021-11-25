create table chainlink_deal(
    deal_id                    bigint        not null,
    deal_cid                   varchar(1000) not null,
    message_cid                varchar(1000) not null,
    height                     bigint        not null,
    piece_cid                  varchar(1000) not null,
    verified_deal              boolean       not null,
    storage_price_per_epoch    varchar(200)  not null,
    signature                  varchar(1000) not null,
    created_at_src             varchar(60)   not null,
    created_at                 bigint        not null, #--precision:second
    piece_size_format          varchar(60)   not null,
    start_height               bigint        not null,
    end_height                 bigint        not null,
    client                     varchar(200)  not null, #--wallet
    client_collateral_format   varchar(60)   not null,
    provider                   varchar(60)   not null,
    provider_tag               varchar(1000),
    verified_provider          int           not null,
    provider_collateral_format varchar(60)   not null,
    status                     int           not null,
    primary key pk_chainlink_deal (deal_id,created_at)
)
partition by range(created_at)(
    partition p20220101 values less than (unix_timestamp("2022-01-01 00:00:00")),
    partition p20230101 values less than (unix_timestamp("2023-01-01 00:00:00")),
    partition p20240101 values less than (unix_timestamp("2024-01-01 00:00:00")),
    partition p20250101 values less than (unix_timestamp("2025-01-01 00:00:00")),
    partition p20260101 values less than (unix_timestamp("2026-01-01 00:00:00"))
)

drop table if exists chainlink_deal;
drop table if exists network;

create table network
(
  id                bigint       not null auto_increment,
  name              varchar(255) not null,
  api_url_prefix    varchar(128) not null,
  description       varchar(2000),
  primary key pk_network (id),
  unique key un_network_api_url_prefix (api_url_prefix)
) engine=InnoDB;

insert into network(name,api_url_prefix) values("calibration", "https://calibration-api.filscout.com/api/v1/storagedeal");

create table chainlink_deal(
    deal_id                    bigint        not null,
    network_id                 bigint        not null,
    deal_cid                   varchar(1000) not null,
    message_cid                varchar(1000) not null,
    height                     bigint        not null,
    piece_cid                  varchar(1000) not null,
    verified_deal              boolean       not null,
    storage_price_per_epoch    varchar(200)  not null,
    signature                  varchar(1000) not null,
    signature_type             varchar(60)   not null,
    created_at_src             varchar(60)   not null,
    created_at                 bigint        not null, #--precision:second
    piece_size_format          varchar(60)   not null,
    start_height               bigint        not null,
    end_height                 bigint        not null,
    client                     varchar(200)  not null, #--wallet
    client_collateral_format   varchar(60)   not null,
    provider                   varchar(60)   not null,
    provider_tag               varchar(1000),
    verified_provider          int           not null,
    provider_collateral_format varchar(60)   not null,
    status                     int           not null,
    primary key pk_chainlink_deal (deal_id),
    constraint fk_chainlink_deal_network_id foreign key (network_id) references network (id)
)

