create index if not exists search_idx on users (first_name varchar_pattern_ops, second_name varchar_pattern_ops);