CREATE OR REPLACE FUNCTION set_updated_at()
  RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- attach to every table that has updated_at
DO $$
BEGIN
FOR tbl IN ARRAY[
      'users','exchanges','exchange_credentials',
      'trading_pairs','order_histories'
    ]
  LOOP
    EXECUTE format(
      'CREATE TRIGGER trg_%1$s_updated_at BEFORE UPDATE ON %1$s FOR EACH ROW EXECUTE FUNCTION set_updated_at()',
      tbl
    );
END LOOP;
END;
$$;
