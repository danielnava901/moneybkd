create or replace function public.get_history_by_code(
    code_input text
    from_date date
    )
returns table (
    country_code text,
    created_at date,
    country_name text,
    value float8
)
language sql
stable
as $$
    select
        h.country_code,
        date_trunc('day', h.created_at)::date as created_at,
        /* Primer country_name del grupo */
        (array_agg(h.country_name order by h.created_at asc))[1] as country_name,
        /* Primer value del grupo */
        (array_agg(h.value order by h.created_at asc))[1] as value
    from history h
    where h.country_code = code_input
        and h.created_at::date >= from_date
    group by
        h.country_code,
        date_trunc('day', h.created_at)
    order by created_at asc;
$$;
