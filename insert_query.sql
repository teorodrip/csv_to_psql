INSERT INTO uncharblog.grid_table
(op, vis, news_heat, quantity, secur, side, status, lmt_pr, tif, fill_qty, avg_pr, filled, working_qty, idle, data_export_restricted, last_update, create_time, vwap, data_export_restricted_2, col_last, bid, ask, volume, d_adv)
SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
WHERE NOT EXISTS (
      SELECT 1 FROM uncharblog.grid_table WHERE op IS NOT DISTINCT FROM $1 AND vis IS NOT DISTINCT FROM $2 AND news_heat IS NOT DISTINCT FROM $3 AND quantity IS NOT DISTINCT FROM $4 AND secur IS NOT DISTINCT FROM $5 AND side IS NOT DISTINCT FROM $6 AND status IS NOT DISTINCT FROM $7 AND lmt_pr IS NOT DISTINCT FROM $8 AND tif IS NOT DISTINCT FROM $9 AND fill_qty IS NOT DISTINCT FROM $10 AND avg_pr IS NOT DISTINCT FROM $11 AND filled IS NOT DISTINCT FROM $12 AND working_qty IS NOT DISTINCT FROM $13 AND idle IS NOT DISTINCT FROM $14 AND data_export_restricted IS NOT DISTINCT FROM $15 AND last_update IS NOT DISTINCT FROM $16 AND create_time IS NOT DISTINCT FROM $17 AND vwap IS NOT DISTINCT FROM $18 AND data_export_restricted_2 IS NOT DISTINCT FROM $19 AND col_last IS NOT DISTINCT FROM $20 AND bid IS NOT DISTINCT FROM $21 AND ask IS NOT DISTINCT FROM $22 AND volume IS NOT DISTINCT FROM $23 AND d_adv IS NOT DISTINCT FROM $24
);
