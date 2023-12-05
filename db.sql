create table tests (
    id int serial primary key,
    name varchar(256) not null,
    date_added timestamp with time zone default current_timestamp,
    constraint unique_name (name),
    check (length(name) > 0)
)


core -> [partitions 
    e.g. A, B,C A -> [list of syn blocks, cma blocks, list of sram blocks, list of rom blocks, list of rf, list of ebb]
    ]

create table fch_block (
    id int serial primary key, 
    name varchar(256) not null,


    constraint unique_name (name),
    check (length(name) > 0)
)




create table syn_block (
    id int serial primary key, 
    name varchar(256) not null,


    constraint unique_name (name),
    check (length(name) > 0)
)


smallest unit 
dcdata0

create table instances (
    id int serial primary key, 
    name varchar(256) not null,
    design_style varchar(256) not null,
    int fch_block_id references fch_blocks(fch_block_id),
    int partition_id references partitions(partition_id),
    


    constraint unique_name (name),
    check (length(name) > 0)
    constraint unique_name (design_style),
    check (length(design_style) > 0)    
)

create table fch_block (
    id int serial primary key, 
    name varchar(256) not null,


    constraint unique_name (name),
    check (length(name) > 0)
)



