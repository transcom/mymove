--
-- PostgreSQL database dump
--

-- Dumped from database version 10.9 (Debian 10.9-1.pgdg90+1)
-- Dumped by pg_dump version 11.5


--
-- Name: stage_conus_to_oconus_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_conus_to_oconus_prices;
CREATE TABLE public.stage_conus_to_oconus_prices (
    origin_domestic_price_area_code text NOT NULL,
    origin_domestic_price_area text NOT NULL,
    destination_intl_price_area_id text NOT NULL,
    destination_intl_price_area text NOT NULL,
    season text NOT NULL,
    hhg_shipping_linehaul_price text NOT NULL,
    ub_price text NOT NULL
);


--
-- Name: stage_counseling_services_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_counseling_services_prices;
CREATE TABLE public.stage_counseling_services_prices (
    contract_year text NOT NULL,
    price_per_task_order text NOT NULL
);


--
-- Name: stage_domestic_international_additional_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_domestic_international_additional_prices;
CREATE TABLE public.stage_domestic_international_additional_prices (
    market text NOT NULL,
    shipment_type text NOT NULL,
    factor text NOT NULL
);


--
-- Name: stage_domestic_linehaul_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_domestic_linehaul_prices;
CREATE TABLE public.stage_domestic_linehaul_prices (
    service_area_number text NOT NULL,
    origin_service_area text NOT NULL,
    services_schedule text NOT NULL,
    season text NOT NULL,
    weight_lower text NOT NULL,
    weight_upper text NOT NULL,
    miles_lower text NOT NULL,
    miles_upper text NOT NULL,
    escalation_number text NOT NULL,
    rate text NOT NULL
);


--
-- Name: stage_domestic_move_accessorial_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_domestic_move_accessorial_prices;
CREATE TABLE public.stage_domestic_move_accessorial_prices (
    services_schedule text NOT NULL,
    service_provided text NOT NULL,
    price_per_unit text NOT NULL
);


--
-- Name: stage_domestic_service_area_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_domestic_service_area_prices;
CREATE TABLE public.stage_domestic_service_area_prices (
    service_area_number text NOT NULL,
    service_area_name text NOT NULL,
    services_schedule text NOT NULL,
    sit_pickup_delivery_schedule text NOT NULL,
    season text NOT NULL,
    shorthaul_price text NOT NULL,
    origin_destination_price text NOT NULL,
    origin_destination_sit_first_day_warehouse text NOT NULL,
    origin_destination_sit_addl_days text NOT NULL
);


--
-- Name: stage_domestic_service_areas; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_domestic_service_areas;
CREATE TABLE public.stage_domestic_service_areas (
    base_point_city text NOT NULL,
    state text NOT NULL,
    service_area_number text NOT NULL,
    zip3s text NOT NULL
);


--
-- Name: stage_domestic_other_pack_prices; Type: TABLE; Schema: public; Owner: -
--
DROP TABLE IF EXISTS public.stage_domestic_other_pack_prices;
CREATE TABLE public.stage_domestic_other_pack_prices (
    services_schedule text NOT NULL,
    service_provided text NOT NULL,
    non_peak_price_per_cwt text NOT NULL,
    peak_price_per_cwt text NOT NULL
);


--
-- Name: stage_domestic_other_sit_prices; Type: TABLE; Schema: public; Owner: -
--
DROP TABLE IF EXISTS public.stage_domestic_other_sit_prices;
CREATE TABLE public.stage_domestic_other_sit_prices (
    sit_pickup_delivery_schedule text NOT NULL,
    service_provided text NOT NULL,
    non_peak_price_per_cwt text NOT NULL,
    peak_price_per_cwt text NOT NULL
);


--
-- Name: stage_international_move_accessorial_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_international_move_accessorial_prices;
CREATE TABLE public.stage_international_move_accessorial_prices (
    market text NOT NULL,
    service_provided text NOT NULL,
    price_per_unit text NOT NULL
);


--
-- Name: stage_international_service_areas; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_international_service_areas;
CREATE TABLE public.stage_international_service_areas (
    rate_area text NOT NULL,
    rate_area_id text NOT NULL
);


--
-- Name: stage_non_standard_locn_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_non_standard_locn_prices;
CREATE TABLE public.stage_non_standard_locn_prices (
    origin_id text NOT NULL,
    origin_area text NOT NULL,
    destination_id text NOT NULL,
    destination_area text NOT NULL,
    move_type text NOT NULL,
    season text NOT NULL,
    hhg_price text NOT NULL,
    ub_price text NOT NULL
);


--
-- Name: stage_oconus_to_conus_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_oconus_to_conus_prices;
CREATE TABLE public.stage_oconus_to_conus_prices (
    origin_intl_price_area_id text NOT NULL,
    origin_intl_price_area text NOT NULL,
    destination_domestic_price_area_area text NOT NULL,
    destination_domestic_price_area text NOT NULL,
    season text NOT NULL,
    hhg_shipping_linehaul_price text NOT NULL,
    ub_price text NOT NULL
);


--
-- Name: stage_oconus_to_oconus_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_oconus_to_oconus_prices;
CREATE TABLE public.stage_oconus_to_oconus_prices (
    origin_intl_price_area_id text NOT NULL,
    origin_intl_price_area text NOT NULL,
    destination_intl_price_area_id text NOT NULL,
    destination_intl_price_area text NOT NULL,
    season text NOT NULL,
    hhg_shipping_linehaul_price text NOT NULL,
    ub_price text NOT NULL
);


--
-- Name: stage_other_intl_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_other_intl_prices;
CREATE TABLE public.stage_other_intl_prices (
    rate_area_code text NOT NULL,
    rate_area_name text NOT NULL,
    hhg_origin_pack_price text NOT NULL,
    hhg_destination_unpack_price text NOT NULL,
    ub_origin_pack_price text NOT NULL,
    ub_destination_unpack_price text NOT NULL,
    origin_destination_sit_first_day_warehouse text NOT NULL,
    origin_destination_sit_addl_days text NOT NULL,
    sit_lte_50_miles text NOT NULL,
    sit_gt_50_miles text NOT NULL,
    season text NOT NULL
);


--
-- Name: stage_price_escalation_discounts; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_price_escalation_discounts;
CREATE TABLE public.stage_price_escalation_discounts (
    contract_year text NOT NULL,
    forecasting_adjustment text NOT NULL,
    discount text NOT NULL,
    price_escalation text NOT NULL
);


--
-- Name: stage_shipment_management_services_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_shipment_management_services_prices;
CREATE TABLE public.stage_shipment_management_services_prices (
    contract_year text NOT NULL,
    price_per_task_order text NOT NULL
);


--
-- Name: stage_transition_prices; Type: TABLE; Schema: public; Owner: -
--

DROP TABLE IF EXISTS public.stage_transition_prices;
CREATE TABLE public.stage_transition_prices (
    contract_year text NOT NULL,
    price_total_cost text NOT NULL
);


--
-- Data for Name: stage_conus_to_oconus_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US16', 'Connecticut', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$16.05', '$44.74');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US16', 'Connecticut', 'AS11', 'New South Wales/Australian Capital Territory', 'Peak', '$18.94', '$52.79');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US16', 'Connecticut', 'GE', 'Germany', 'NonPeak', '$24.38', '$31.40');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US16', 'Connecticut', 'GE', 'Germany', 'Peak', '$28.77', '$37.05');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US16', 'Connecticut', 'US8101000', 'Alaska (Zone) I', 'NonPeak', '$27.63', '$34.11');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US16', 'Connecticut', 'US8101000', 'Alaska (Zone) I', 'Peak', '$32.60', '$40.25');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US47', 'Alabama', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$30.90', '$33.98');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US47', 'Alabama', 'AS11', 'New South Wales/Australian Capital Territory', 'Peak', '$36.46', '$40.10');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US47', 'Alabama', 'GE', 'Germany', 'NonPeak', '$17.57', '$44.91');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US47', 'Alabama', 'GE', 'Germany', 'Peak', '$20.73', '$52.99');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US47', 'Alabama', 'US8101000', 'Alaska (Zone) I', 'NonPeak', '$24.38', '$32.44');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US47', 'Alabama', 'US8101000', 'Alaska (Zone) I', 'Peak', '$28.77', '$38.28');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US4965500', 'Florida Keys', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$27.63', '$32.44');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US4965500', 'Florida Keys', 'AS11', 'New South Wales/Australian Capital Territory', 'Peak', '$32.60', '$38.28');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US4965500', 'Florida Keys', 'GE', 'Germany', 'NonPeak', '$16.05', '$34.33');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US4965500', 'Florida Keys', 'GE', 'Germany', 'Peak', '$18.94', '$40.51');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US4965500', 'Florida Keys', 'US8101000', 'Alaska (Zone) I', 'NonPeak', '$17.57', '$44.74');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US4965500', 'Florida Keys', 'US8101000', 'Alaska (Zone) I', 'Peak', '$20.73', '$52.79');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US68', 'Texas-South', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$16.05', '$31.29');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US68', 'Texas-South', 'AS11', 'New South Wales/Australian Capital Territory', 'Peak', '$18.94', '$36.92');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US68', 'Texas-South', 'GE', 'Germany', 'NonPeak', '$24.38', '$34.45');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US68', 'Texas-South', 'GE', 'Germany', 'Peak', '$28.77', '$40.65');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US68', 'Texas-South', 'US8101000', 'Alaska (Zone) I', 'NonPeak', '$27.63', '$44.91');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US68', 'Texas-South', 'US8101000', 'Alaska (Zone) I', 'Peak', '$32.60', '$52.99');

-- Data for Test_mapZipCodesToReRateAreas
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US34', 'Ohio', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$16.05', '$44.74');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US49', 'Florida', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$27.63', '$32.44');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US4964400', 'Florida', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$27.63', '$32.44');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US51', 'North Dakota', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$30.90', '$33.98');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US56', 'Missouri', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$27.63', '$32.44');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US66', 'Texas', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$27.63', '$32.44');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US87', 'California', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$27.63', '$32.44');
INSERT INTO public.stage_conus_to_oconus_prices (origin_domestic_price_area_code, origin_domestic_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US88', 'California', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$27.63', '$32.44');

--
-- Data for Name: stage_counseling_services_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_counseling_services_prices (contract_year, price_per_task_order) VALUES ('Base Period Year 1', '222.63');
INSERT INTO public.stage_counseling_services_prices (contract_year, price_per_task_order) VALUES ('Base Period Year 2', '223.53');
INSERT INTO public.stage_counseling_services_prices (contract_year, price_per_task_order) VALUES ('Base Period Year 3', '224.69');
INSERT INTO public.stage_counseling_services_prices (contract_year, price_per_task_order) VALUES ('Option Period 1', '224.88');
INSERT INTO public.stage_counseling_services_prices (contract_year, price_per_task_order) VALUES ('Option Period 2', '226.77');
INSERT INTO public.stage_counseling_services_prices (contract_year, price_per_task_order) VALUES ('Award Term 1', '227.59');
INSERT INTO public.stage_counseling_services_prices (contract_year, price_per_task_order) VALUES ('Award Term 2', '227.99');
INSERT INTO public.stage_counseling_services_prices (contract_year, price_per_task_order) VALUES ('Option Period 3', '228.17');


--
-- Data for Name: stage_domestic_international_additional_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_domestic_international_additional_prices (market, shipment_type, factor) VALUES ('CONUS', 'Mobile Homes', '1.20');
INSERT INTO public.stage_domestic_international_additional_prices (market, shipment_type, factor) VALUES ('CONUS', 'Tow Away Boat Service', '1.10');
INSERT INTO public.stage_domestic_international_additional_prices (market, shipment_type, factor) VALUES ('CONUS', 'Haul Away Boat Service', '1.30');
INSERT INTO public.stage_domestic_international_additional_prices (market, shipment_type, factor) VALUES ('OCONUS', 'Tow Away Boat Service', '1.32');
INSERT INTO public.stage_domestic_international_additional_prices (market, shipment_type, factor) VALUES ('OCONUS', 'Haul Away Boat Service', '1.40');
INSERT INTO public.stage_domestic_international_additional_prices (market, shipment_type, factor) VALUES ('CONUS', 'NTS Packing Factor', '1.32');
INSERT INTO public.stage_domestic_international_additional_prices (market, shipment_type, factor) VALUES ('OCONUS', 'NTS Packing Factor', '1.45');


--
-- Data for Name: stage_domestic_linehaul_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '500', '4999', '0', '250', '0', '$2.477');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '500', '4999', '251', '500', '0', '$2.727');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '500', '4999', '501', '1000', '0', '$3.228');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '500', '4999', '1001', '1500', '0', '$3.728');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '500', '4999', '1501', '2000', '0', '$4.228');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '500', '4999', '2001', '2500', '0', '$4.728');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '500', '4999', '2501', '3000', '0', '$5.228');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '500', '4999', '3001', '3500', '0', '$5.728');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '500', '4999', '3501', '4000', '0', '$6.228');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '500', '4999', '4001', '999999', '0', '$27.132');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '5000', '9999', '0', '250', '0', '$4.705');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '5000', '9999', '251', '500', '0', '$4.955');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '5000', '9999', '501', '1000', '0', '$5.455');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '5000', '9999', '1001', '1500', '0', '$5.955');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '5000', '9999', '1501', '2000', '0', '$6.455');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '5000', '9999', '2001', '2500', '0', '$6.956');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '5000', '9999', '2501', '3000', '0', '$7.456');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '5000', '9999', '3001', '3500', '0', '$7.956');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '5000', '9999', '3501', '4000', '0', '$8.456');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '5000', '9999', '4001', '999999', '0', '$29.360');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '10000', '999999', '0', '250', '0', '$11.389');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '10000', '999999', '251', '500', '0', '$11.639');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '10000', '999999', '501', '1000', '0', '$12.139');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '10000', '999999', '1001', '1500', '0', '$12.639');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '10000', '999999', '1501', '2000', '0', '$13.139');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '10000', '999999', '2001', '2500', '0', '$13.639');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '10000', '999999', '2501', '3000', '0', '$14.139');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '10000', '999999', '3001', '3500', '0', '$14.640');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '10000', '999999', '3501', '4000', '0', '$15.140');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'NonPeak', '10000', '999999', '4001', '999999', '0', '$36.044');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '500', '4999', '0', '250', '0', '$2.583');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '500', '4999', '251', '500', '0', '$2.843');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '500', '4999', '501', '1000', '0', '$3.365');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '500', '4999', '1001', '1500', '0', '$3.886');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '500', '4999', '1501', '2000', '0', '$4.408');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '500', '4999', '2001', '2500', '0', '$4.929');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '500', '4999', '2501', '3000', '0', '$5.450');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '500', '4999', '3001', '3500', '0', '$5.972');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '500', '4999', '3501', '4000', '0', '$6.493');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '500', '4999', '4001', '999999', '0', '$28.287');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '5000', '9999', '0', '250', '0', '$4.905');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '5000', '9999', '251', '500', '0', '$5.166');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '5000', '9999', '501', '1000', '0', '$5.687');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '5000', '9999', '1001', '1500', '0', '$6.209');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '5000', '9999', '1501', '2000', '0', '$6.730');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '5000', '9999', '2001', '2500', '0', '$7.252');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '5000', '9999', '2501', '3000', '0', '$7.773');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '5000', '9999', '3001', '3500', '0', '$8.294');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '5000', '9999', '3501', '4000', '0', '$8.816');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '5000', '9999', '4001', '999999', '0', '$30.610');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '10000', '999999', '0', '250', '0', '$11.874');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '10000', '999999', '251', '500', '0', '$12.134');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '10000', '999999', '501', '1000', '0', '$12.656');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '10000', '999999', '1001', '1500', '0', '$13.177');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '10000', '999999', '1501', '2000', '0', '$13.698');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '10000', '999999', '2001', '2500', '0', '$14.220');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '10000', '999999', '2501', '3000', '0', '$14.741');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '10000', '999999', '3001', '3500', '0', '$15.263');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '10000', '999999', '3501', '4000', '0', '$15.784');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('4.0', 'Birmingham, AL', '2', 'Peak', '10000', '999999', '4001', '999999', '0', '$37.578');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '500', '4999', '0', '250', '0', '$2.477');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '500', '4999', '251', '500', '0', '$2.727');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '500', '4999', '501', '1000', '0', '$3.228');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '500', '4999', '1001', '1500', '0', '$3.728');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '500', '4999', '1501', '2000', '0', '$4.228');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '500', '4999', '2001', '2500', '0', '$4.728');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '500', '4999', '2501', '3000', '0', '$5.228');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '500', '4999', '3001', '3500', '0', '$5.728');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '500', '4999', '3501', '4000', '0', '$6.228');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '500', '4999', '4001', '999999', '0', '$27.132');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '5000', '9999', '0', '250', '0', '$4.705');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '5000', '9999', '251', '500', '0', '$4.955');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '5000', '9999', '501', '1000', '0', '$5.455');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '5000', '9999', '1001', '1500', '0', '$5.955');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '5000', '9999', '1501', '2000', '0', '$6.455');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '5000', '9999', '2001', '2500', '0', '$6.956');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '5000', '9999', '2501', '3000', '0', '$7.456');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '5000', '9999', '3001', '3500', '0', '$7.956');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '5000', '9999', '3501', '4000', '0', '$8.456');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '5000', '9999', '4001', '999999', '0', '$29.360');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '10000', '999999', '0', '250', '0', '$11.389');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '10000', '999999', '251', '500', '0', '$11.639');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '10000', '999999', '501', '1000', '0', '$12.139');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '10000', '999999', '1001', '1500', '0', '$12.639');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '10000', '999999', '1501', '2000', '0', '$13.139');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '10000', '999999', '2001', '2500', '0', '$13.639');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '10000', '999999', '2501', '3000', '0', '$14.139');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '10000', '999999', '3001', '3500', '0', '$14.640');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '10000', '999999', '3501', '4000', '0', '$15.140');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'NonPeak', '10000', '999999', '4001', '999999', '0', '$36.044');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '500', '4999', '0', '250', '0', '$2.583');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '500', '4999', '251', '500', '0', '$2.843');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '500', '4999', '501', '1000', '0', '$3.365');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '500', '4999', '1001', '1500', '0', '$3.886');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '500', '4999', '1501', '2000', '0', '$4.408');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '500', '4999', '2001', '2500', '0', '$4.929');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '500', '4999', '2501', '3000', '0', '$5.450');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '500', '4999', '3001', '3500', '0', '$5.972');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '500', '4999', '3501', '4000', '0', '$6.493');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '500', '4999', '4001', '999999', '0', '$28.287');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '5000', '9999', '0', '250', '0', '$4.905');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '5000', '9999', '251', '500', '0', '$5.166');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '5000', '9999', '501', '1000', '0', '$5.687');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '5000', '9999', '1001', '1500', '0', '$6.209');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '5000', '9999', '1501', '2000', '0', '$6.730');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '5000', '9999', '2001', '2500', '0', '$7.252');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '5000', '9999', '2501', '3000', '0', '$7.773');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '5000', '9999', '3001', '3500', '0', '$8.294');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '5000', '9999', '3501', '4000', '0', '$8.816');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '5000', '9999', '4001', '999999', '0', '$30.610');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '10000', '999999', '0', '250', '0', '$11.874');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '10000', '999999', '251', '500', '0', '$12.134');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '10000', '999999', '501', '1000', '0', '$12.656');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '10000', '999999', '1001', '1500', '0', '$13.177');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '10000', '999999', '1501', '2000', '0', '$13.698');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '10000', '999999', '2001', '2500', '0', '$14.220');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '10000', '999999', '2501', '3000', '0', '$14.741');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '10000', '999999', '3001', '3500', '0', '$15.263');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '10000', '999999', '3501', '4000', '0', '$15.784');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('452.0', 'Springfield, MO', '1', 'Peak', '10000', '999999', '4001', '999999', '0', '$37.578');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '500', '4999', '0', '250', '0', '$2.161');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '500', '4999', '251', '500', '0', '$2.379');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '500', '4999', '501', '1000', '0', '$2.815');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '500', '4999', '1001', '1500', '0', '$3.252');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '500', '4999', '1501', '2000', '0', '$3.688');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '500', '4999', '2001', '2500', '0', '$4.124');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '500', '4999', '2501', '3000', '0', '$4.560');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '500', '4999', '3001', '3500', '0', '$4.997');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '500', '4999', '3501', '4000', '0', '$5.433');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '500', '4999', '4001', '999999', '0', '$23.669');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '5000', '9999', '0', '250', '0', '$4.105');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '5000', '9999', '251', '500', '0', '$4.323');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '5000', '9999', '501', '1000', '0', '$4.759');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '5000', '9999', '1001', '1500', '0', '$5.195');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '5000', '9999', '1501', '2000', '0', '$5.631');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '5000', '9999', '2001', '2500', '0', '$6.068');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '5000', '9999', '2501', '3000', '0', '$6.504');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '5000', '9999', '3001', '3500', '0', '$6.940');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '5000', '9999', '3501', '4000', '0', '$7.376');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '5000', '9999', '4001', '999999', '0', '$25.612');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '10000', '999999', '0', '250', '0', '$9.935');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '10000', '999999', '251', '500', '0', '$10.153');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '10000', '999999', '501', '1000', '0', '$10.589');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '10000', '999999', '1001', '1500', '0', '$11.026');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '10000', '999999', '1501', '2000', '0', '$11.462');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '10000', '999999', '2001', '2500', '0', '$11.898');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '10000', '999999', '2501', '3000', '0', '$12.334');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '10000', '999999', '3001', '3500', '0', '$12.771');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '10000', '999999', '3501', '4000', '0', '$13.207');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'NonPeak', '10000', '999999', '4001', '999999', '0', '$31.443');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '500', '4999', '0', '250', '0', '$2.293');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '500', '4999', '251', '500', '0', '$2.524');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '500', '4999', '501', '1000', '0', '$2.987');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '500', '4999', '1001', '1500', '0', '$3.450');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '500', '4999', '1501', '2000', '0', '$3.913');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '500', '4999', '2001', '2500', '0', '$4.376');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '500', '4999', '2501', '3000', '0', '$4.839');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '500', '4999', '3001', '3500', '0', '$5.301');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '500', '4999', '3501', '4000', '0', '$5.764');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '500', '4999', '4001', '999999', '0', '$25.112');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '5000', '9999', '0', '250', '0', '$4.355');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '5000', '9999', '251', '500', '0', '$4.586');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '5000', '9999', '501', '1000', '0', '$5.049');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '5000', '9999', '1001', '1500', '0', '$5.512');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '5000', '9999', '1501', '2000', '0', '$5.975');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '5000', '9999', '2001', '2500', '0', '$6.438');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '5000', '9999', '2501', '3000', '0', '$6.900');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '5000', '9999', '3001', '3500', '0', '$7.363');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '5000', '9999', '3501', '4000', '0', '$7.826');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '5000', '9999', '4001', '999999', '0', '$27.174');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '10000', '999999', '0', '250', '0', '$10.541');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '10000', '999999', '251', '500', '0', '$10.772');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '10000', '999999', '501', '1000', '0', '$11.235');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '10000', '999999', '1001', '1500', '0', '$11.698');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '10000', '999999', '1501', '2000', '0', '$12.161');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '10000', '999999', '2001', '2500', '0', '$12.624');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '10000', '999999', '2501', '3000', '0', '$13.087');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '10000', '999999', '3001', '3500', '0', '$13.549');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '10000', '999999', '3501', '4000', '0', '$14.012');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('592.0', 'Dickinson, ND', '3', 'Peak', '10000', '999999', '4001', '999999', '0', '$33.360');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '500', '4999', '0', '250', '0', '$2.161');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '500', '4999', '251', '500', '0', '$2.379');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '500', '4999', '501', '1000', '0', '$2.815');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '500', '4999', '1001', '1500', '0', '$3.252');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '500', '4999', '1501', '2000', '0', '$3.688');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '500', '4999', '2001', '2500', '0', '$4.124');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '500', '4999', '2501', '3000', '0', '$4.560');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '500', '4999', '3001', '3500', '0', '$4.997');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '500', '4999', '3501', '4000', '0', '$5.433');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '500', '4999', '4001', '999999', '0', '$23.669');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '5000', '9999', '0', '250', '0', '$4.105');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '5000', '9999', '251', '500', '0', '$4.323');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '5000', '9999', '501', '1000', '0', '$4.759');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '5000', '9999', '1001', '1500', '0', '$5.195');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '5000', '9999', '1501', '2000', '0', '$5.631');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '5000', '9999', '2001', '2500', '0', '$6.068');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '5000', '9999', '2501', '3000', '0', '$6.504');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '5000', '9999', '3001', '3500', '0', '$6.940');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '5000', '9999', '3501', '4000', '0', '$7.376');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '5000', '9999', '4001', '999999', '0', '$25.612');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '10000', '999999', '0', '250', '0', '$9.935');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '10000', '999999', '251', '500', '0', '$10.153');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '10000', '999999', '501', '1000', '0', '$10.589');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '10000', '999999', '1001', '1500', '0', '$11.026');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '10000', '999999', '1501', '2000', '0', '$11.462');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '10000', '999999', '2001', '2500', '0', '$11.898');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '10000', '999999', '2501', '3000', '0', '$12.334');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '10000', '999999', '3001', '3500', '0', '$12.771');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '10000', '999999', '3501', '4000', '0', '$13.207');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'NonPeak', '10000', '999999', '4001', '999999', '0', '$31.443');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '500', '4999', '0', '250', '0', '$2.293');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '500', '4999', '251', '500', '0', '$2.524');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '500', '4999', '501', '1000', '0', '$2.987');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '500', '4999', '1001', '1500', '0', '$3.450');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '500', '4999', '1501', '2000', '0', '$3.913');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '500', '4999', '2001', '2500', '0', '$4.376');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '500', '4999', '2501', '3000', '0', '$4.839');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '500', '4999', '3001', '3500', '0', '$5.301');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '500', '4999', '3501', '4000', '0', '$5.764');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '500', '4999', '4001', '999999', '0', '$25.112');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '5000', '9999', '0', '250', '0', '$4.355');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '5000', '9999', '251', '500', '0', '$4.586');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '5000', '9999', '501', '1000', '0', '$5.049');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '5000', '9999', '1001', '1500', '0', '$5.512');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '5000', '9999', '1501', '2000', '0', '$5.975');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '5000', '9999', '2001', '2500', '0', '$6.438');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '5000', '9999', '2501', '3000', '0', '$6.900');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '5000', '9999', '3001', '3500', '0', '$7.363');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '5000', '9999', '3501', '4000', '0', '$7.826');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '5000', '9999', '4001', '999999', '0', '$27.174');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '10000', '999999', '0', '250', '0', '$10.541');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '10000', '999999', '251', '500', '0', '$10.772');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '10000', '999999', '501', '1000', '0', '$11.235');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '10000', '999999', '1001', '1500', '0', '$11.698');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '10000', '999999', '1501', '2000', '0', '$12.161');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '10000', '999999', '2001', '2500', '0', '$12.624');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '10000', '999999', '2501', '3000', '0', '$13.087');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '10000', '999999', '3001', '3500', '0', '$13.549');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '10000', '999999', '3501', '4000', '0', '$14.012');
INSERT INTO public.stage_domestic_linehaul_prices (service_area_number, origin_service_area, services_schedule, season, weight_lower, weight_upper, miles_lower, miles_upper, escalation_number, rate) VALUES ('616.0', 'Columbus, OH', '2', 'Peak', '10000', '999999', '4001', '999999', '0', '$33.360');


--
-- Data for Name: stage_domestic_move_accessorial_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_domestic_move_accessorial_prices (services_schedule, service_provided, price_per_unit) VALUES ('1', 'Crating (per cubic ft.)', '23.69');
INSERT INTO public.stage_domestic_move_accessorial_prices (services_schedule, service_provided, price_per_unit) VALUES ('2', 'Crating (per cubic ft.)', '23.69');
INSERT INTO public.stage_domestic_move_accessorial_prices (services_schedule, service_provided, price_per_unit) VALUES ('3', 'Crating (per cubic ft.)', '23.69');
INSERT INTO public.stage_domestic_move_accessorial_prices (services_schedule, service_provided, price_per_unit) VALUES ('1', 'Uncrating (per cubic ft.)', '5.95');
INSERT INTO public.stage_domestic_move_accessorial_prices (services_schedule, service_provided, price_per_unit) VALUES ('2', 'Uncrating (per cubic ft.)', '5.95');
INSERT INTO public.stage_domestic_move_accessorial_prices (services_schedule, service_provided, price_per_unit) VALUES ('3', 'Uncrating (per cubic ft.)', '5.95');
INSERT INTO public.stage_domestic_move_accessorial_prices (services_schedule, service_provided, price_per_unit) VALUES ('1', 'Shuttle Service (per cwt)', '5.05');
INSERT INTO public.stage_domestic_move_accessorial_prices (services_schedule, service_provided, price_per_unit) VALUES ('2', 'Shuttle Service (per cwt)', '5.41');
INSERT INTO public.stage_domestic_move_accessorial_prices (services_schedule, service_provided, price_per_unit) VALUES ('3', 'Shuttle Service (per cwt)', '5.76');


--
-- Data for Name: stage_domestic_service_area_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_domestic_service_area_prices (service_area_number, service_area_name, services_schedule, sit_pickup_delivery_schedule, season, shorthaul_price, origin_destination_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days) VALUES ('184.0', 'Sanford, FL', '2', '2', 'NonPeak', '$1.27', '$6.89', '$19.31', '$0.68');
INSERT INTO public.stage_domestic_service_area_prices (service_area_number, service_area_name, services_schedule, sit_pickup_delivery_schedule, season, shorthaul_price, origin_destination_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days) VALUES ('184.0', 'Sanford, FL', '2', '2', 'Peak', '$1.28', '$8.28', '$22.31', '$0.75');
INSERT INTO public.stage_domestic_service_area_prices (service_area_number, service_area_name, services_schedule, sit_pickup_delivery_schedule, season, shorthaul_price, origin_destination_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days) VALUES ('4.0', 'Birmingham, AL', '2', '2', 'NonPeak', '$1.27', '$6.89', '$19.31', '$0.68');
INSERT INTO public.stage_domestic_service_area_prices (service_area_number, service_area_name, services_schedule, sit_pickup_delivery_schedule, season, shorthaul_price, origin_destination_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days) VALUES ('4.0', 'Birmingham, AL', '2', '2', 'Peak', '$1.46', '$7.92', '$22.21', '$0.78');
INSERT INTO public.stage_domestic_service_area_prices (service_area_number, service_area_name, services_schedule, sit_pickup_delivery_schedule, season, shorthaul_price, origin_destination_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days) VALUES ('452.0', 'Springfield, MO', '1', '3', 'NonPeak', '$1.08', '$7.20', '$14.27', '$0.55');
INSERT INTO public.stage_domestic_service_area_prices (service_area_number, service_area_name, services_schedule, sit_pickup_delivery_schedule, season, shorthaul_price, origin_destination_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days) VALUES ('452.0', 'Springfield, MO', '1', '3', 'Peak', '$1.24', '$8.28', '$16.41', '$0.63');
INSERT INTO public.stage_domestic_service_area_prices (service_area_number, service_area_name, services_schedule, sit_pickup_delivery_schedule, season, shorthaul_price, origin_destination_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days) VALUES ('592.0', 'Dickinson, ND', '3', '3', 'NonPeak', '$0.16', '$5.81', '$15.97', '$0.62');
INSERT INTO public.stage_domestic_service_area_prices (service_area_number, service_area_name, services_schedule, sit_pickup_delivery_schedule, season, shorthaul_price, origin_destination_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days) VALUES ('592.0', 'Dickinson, ND', '3', '3', 'Peak', '$0.18', '$6.68', '$18.37', '$0.71');
INSERT INTO public.stage_domestic_service_area_prices (service_area_number, service_area_name, services_schedule, sit_pickup_delivery_schedule, season, shorthaul_price, origin_destination_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days) VALUES ('616.0', 'Columbus, OH', '2', '2', 'NonPeak', '$2.44', '$9.55', '$15.81', '$0.62');
INSERT INTO public.stage_domestic_service_area_prices (service_area_number, service_area_name, services_schedule, sit_pickup_delivery_schedule, season, shorthaul_price, origin_destination_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days) VALUES ('616.0', 'Columbus, OH', '2', '2', 'Peak', '$2.81', '$10.98', '$18.18', '$0.71');


--
-- Data for Name: stage_domestic_service_areas; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Sanford', 'FL', '184', '327.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Birmingham', 'AL', '004', '352.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Butler', 'MO', '452.0', '647.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Carbon Hill', 'AL', '004', '355.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Chulafinnee', 'AL', '004', '362.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Collinsville', 'AL', '004', '359.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Columbus', 'OH', '616.0', '432.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Cumberland', 'OH', '616.0', '437.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Dickinson', 'ND', '592.0', '586.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('La Rue', 'OH', '616.0', '433.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Leeds', 'AL', '004', '350,351');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Neosho', 'MO', '452.0', '648.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('South Bloomfield', 'OH', '616.0', '431.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Springfield', 'MO', '452.0', '656,657,658');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Tuscaloosa', 'AL', '004', '354.0');
INSERT INTO public.stage_domestic_service_areas (base_point_city, state, service_area_number, zip3s) VALUES ('Worthington', 'OH', '616.0', '430.0');


--
-- Data for Name: stage_domestic_other_pack_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_domestic_other_pack_prices (services_schedule, service_provided, non_peak_price_per_cwt, peak_price_per_cwt) VALUES ('1','Packing (per cwt)','$63.33','$65.44');
INSERT INTO public.stage_domestic_other_pack_prices (services_schedule, service_provided, non_peak_price_per_cwt, peak_price_per_cwt) VALUES ('2','Packing (per cwt)','$72.50','$73.20');
INSERT INTO public.stage_domestic_other_pack_prices (services_schedule, service_provided, non_peak_price_per_cwt, peak_price_per_cwt) VALUES ('3','Packing (per cwt)','$73.95','$80.00');
INSERT INTO public.stage_domestic_other_pack_prices (services_schedule, service_provided, non_peak_price_per_cwt, peak_price_per_cwt) VALUES ('1','Unpack (per cwt)','$83.34','$85.44');
INSERT INTO public.stage_domestic_other_pack_prices (services_schedule, service_provided, non_peak_price_per_cwt, peak_price_per_cwt) VALUES ('2','Unpack (per cwt)','$5.97','$6.50');
INSERT INTO public.stage_domestic_other_pack_prices (services_schedule, service_provided, non_peak_price_per_cwt, peak_price_per_cwt) VALUES ('3','Unpack (per cwt)','$5.97','$6.50');


--
-- Data for Name: stage_domestic_other_sit_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_domestic_other_sit_prices (sit_pickup_delivery_schedule, service_provided, non_peak_price_per_cwt, peak_price_per_cwt) VALUES ('1','SIT Pickup / Delivery ≤50 miles (per cwt)','$217.96','$220.11');
INSERT INTO public.stage_domestic_other_sit_prices (sit_pickup_delivery_schedule, service_provided, non_peak_price_per_cwt, peak_price_per_cwt) VALUES ('2','SIT Pickup / Delivery ≤50 miles (per cwt)','$234.40','$241.22');
INSERT INTO public.stage_domestic_other_sit_prices (sit_pickup_delivery_schedule, service_provided, non_peak_price_per_cwt, peak_price_per_cwt) VALUES ('3','SIT Pickup / Delivery ≤50 miles (per cwt)','$246.25','$250.30');


--
-- Data for Name: stage_international_move_accessorial_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_international_move_accessorial_prices (market, service_provided, price_per_unit) VALUES ('CONUS', 'Crating (per cubic ft.)', '25.61');
INSERT INTO public.stage_international_move_accessorial_prices (market, service_provided, price_per_unit) VALUES ('OCONUS', 'Crating (per cubic ft.)', '28.59');
INSERT INTO public.stage_international_move_accessorial_prices (market, service_provided, price_per_unit) VALUES ('CONUS', 'Uncrating (per cubic ft.)', '6.54');
INSERT INTO public.stage_international_move_accessorial_prices (market, service_provided, price_per_unit) VALUES ('OCONUS', 'Uncrating (per cubic ft.)', '6.54');
INSERT INTO public.stage_international_move_accessorial_prices (market, service_provided, price_per_unit) VALUES ('CONUS', 'Shuttle Service (per cwt)', '145.29');
INSERT INTO public.stage_international_move_accessorial_prices (market, service_provided, price_per_unit) VALUES ('OCONUS', 'Shuttle Service (per cwt)', '156.23');


--
-- Data for Name: stage_international_service_areas; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_international_service_areas (rate_area, rate_area_id) VALUES ('Alaska (Zone) I', 'US8101000');
INSERT INTO public.stage_international_service_areas (rate_area, rate_area_id) VALUES ('Germany', 'GE');
INSERT INTO public.stage_international_service_areas (rate_area, rate_area_id) VALUES ('New South Wales/Australian Capital Territory', 'AS11');
INSERT INTO public.stage_international_service_areas (rate_area, rate_area_id) VALUES ('Canada Central', 'NSRA2');
INSERT INTO public.stage_international_service_areas (rate_area, rate_area_id) VALUES ('Pacific Islands', 'NSRA13');


--
-- Data for Name: stage_non_standard_locn_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'NSRA2', 'Canada Central', 'NSRA to NSRA', 'NonPeak', '$9.77', '$44.58');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'NSRA2', 'Canada Central', 'NSRA to NSRA', 'Peak', '$11.53', '$52.60');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'NSRA13', 'Pacific Islands', 'NSRA to NSRA', 'NonPeak', '$48.49', '$47.93');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'NSRA13', 'Pacific Islands', 'NSRA to NSRA', 'Peak', '$57.22', '$56.56');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'NSRA2', 'Canada Central', 'NSRA to NSRA', 'NonPeak', '$9.77', '$11.66');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'NSRA2', 'Canada Central', 'NSRA to NSRA', 'Peak', '$11.53', '$13.76');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'NSRA13', 'Pacific Islands', 'NSRA to NSRA', 'NonPeak', '$48.49', '$12.03');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'NSRA13', 'Pacific Islands', 'NSRA to NSRA', 'Peak', '$57.22', '$14.20');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'AS11', 'New South Wales/Australian Capital Territory', 'NSRA to OCONUS', 'NonPeak', '$51.21', '$61.80');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'AS11', 'New South Wales/Australian Capital Territory', 'NSRA to OCONUS', 'Peak', '$60.43', '$72.92');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'GE', 'Germany', 'NSRA to OCONUS', 'NonPeak', '$16.05', '$31.29');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'GE', 'Germany', 'NSRA to OCONUS', 'Peak', '$18.94', '$36.92');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'US8101000', 'Alaska (Zone) I', 'NSRA to OCONUS', 'NonPeak', '$15.28', '$47.93');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'US8101000', 'Alaska (Zone) I', 'NSRA to OCONUS', 'Peak', '$18.03', '$56.56');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'AS11', 'New South Wales/Australian Capital Territory', 'NSRA to OCONUS', 'NonPeak', '$51.72', '$11.75');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'AS11', 'New South Wales/Australian Capital Territory', 'NSRA to OCONUS', 'Peak', '$61.03', '$13.86');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'GE', 'Germany', 'NSRA to OCONUS', 'NonPeak', '$17.57', '$47.76');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'GE', 'Germany', 'NSRA to OCONUS', 'Peak', '$20.73', '$56.36');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'US8101000', 'Alaska (Zone) I', 'NSRA to OCONUS', 'NonPeak', '$16.80', '$61.92');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'US8101000', 'Alaska (Zone) I', 'NSRA to OCONUS', 'Peak', '$19.82', '$73.07');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'NSRA2', 'Canada Central', 'OCONUS to NSRA', 'NonPeak', '$9.77', '$44.58');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'NSRA2', 'Canada Central', 'OCONUS to NSRA', 'Peak', '$11.53', '$52.60');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'NSRA13', 'Pacific Islands', 'OCONUS to NSRA', 'NonPeak', '$48.49', '$47.93');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'NSRA13', 'Pacific Islands', 'OCONUS to NSRA', 'Peak', '$57.22', '$56.56');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('GE', 'Germany', 'NSRA2', 'Canada Central', 'OCONUS to NSRA', 'NonPeak', '$48.72', '$10.50');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('GE', 'Germany', 'NSRA2', 'Canada Central', 'OCONUS to NSRA', 'Peak', '$57.49', '$12.39');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('GE', 'Germany', 'NSRA13', 'Pacific Islands', 'OCONUS to NSRA', 'NonPeak', '$16.00', '$47.93');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('GE', 'Germany', 'NSRA13', 'Pacific Islands', 'OCONUS to NSRA', 'Peak', '$18.88', '$56.56');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'NSRA2', 'Canada Central', 'OCONUS to NSRA', 'NonPeak', '$51.01', '$44.58');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'NSRA2', 'Canada Central', 'OCONUS to NSRA', 'Peak', '$60.19', '$52.60');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'NSRA13', 'Pacific Islands', 'OCONUS to NSRA', 'NonPeak', '$71.61', '$47.93');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'NSRA13', 'Pacific Islands', 'OCONUS to NSRA', 'Peak', '$84.50', '$56.56');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'US16', 'Connecticut', 'NSRA to CONUS', 'NonPeak', '$47.42', '$44.58');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'US16', 'Connecticut', 'NSRA to CONUS', 'Peak', '$55.96', '$52.60');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'US47', 'Alabama', 'NSRA to CONUS', 'NonPeak', '$10.83', '$10.25');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'US47', 'Alabama', 'NSRA to CONUS', 'Peak', '$12.78', '$12.09');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'US4965500', 'Florida Keys', 'NSRA to CONUS', 'NonPeak', '$9.31', '$17.17');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'US4965500', 'Florida Keys', 'NSRA to CONUS', 'Peak', '$10.99', '$20.26');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'US68', 'Texas-South', 'NSRA to CONUS', 'NonPeak', '$48.49', '$47.93');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA2', 'Canada Central', 'US68', 'Texas-South', 'NSRA to CONUS', 'Peak', '$57.22', '$56.56');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'US16', 'Connecticut', 'NSRA to CONUS', 'NonPeak', '$8.88', '$43.94');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'US16', 'Connecticut', 'NSRA to CONUS', 'Peak', '$10.48', '$51.85');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'US47', 'Alabama', 'NSRA to CONUS', 'NonPeak', '$15.28', '$12.03');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'US47', 'Alabama', 'NSRA to CONUS', 'Peak', '$18.03', '$14.20');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'US4965500', 'Florida Keys', 'NSRA to CONUS', 'NonPeak', '$18.33', '$44.58');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'US4965500', 'Florida Keys', 'NSRA to CONUS', 'Peak', '$21.63', '$52.60');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'US68', 'Texas-South', 'NSRA to CONUS', 'NonPeak', '$9.77', '$44.58');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('NSRA13', 'Pacific Islands', 'US68', 'Texas-South', 'NSRA to CONUS', 'Peak', '$11.53', '$52.60');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US16', 'Connecticut', 'NSRA2', 'Canada Central', 'CONUS to NSRA', 'NonPeak', '$10.35', '$11.75');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US16', 'Connecticut', 'NSRA2', 'Canada Central', 'CONUS to NSRA', 'Peak', '$12.21', '$13.86');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US16', 'Connecticut', 'NSRA13', 'Pacific Islands', 'CONUS to NSRA', 'NonPeak', '$51.55', '$47.93');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US16', 'Connecticut', 'NSRA13', 'Pacific Islands', 'CONUS to NSRA', 'Peak', '$60.83', '$56.56');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US47', 'Alabama', 'NSRA2', 'Canada Central', 'CONUS to NSRA', 'NonPeak', '$71.89', '$11.66');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US47', 'Alabama', 'NSRA2', 'Canada Central', 'CONUS to NSRA', 'Peak', '$84.83', '$13.76');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US47', 'Alabama', 'NSRA13', 'Pacific Islands', 'CONUS to NSRA', 'NonPeak', '$30.90', '$12.03');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US47', 'Alabama', 'NSRA13', 'Pacific Islands', 'CONUS to NSRA', 'Peak', '$36.46', '$14.20');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US4965500', 'Florida Keys', 'NSRA2', 'Canada Central', 'CONUS to NSRA', 'NonPeak', '$47.42', '$44.24');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US4965500', 'Florida Keys', 'NSRA2', 'Canada Central', 'CONUS to NSRA', 'Peak', '$55.96', '$52.20');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US4965500', 'Florida Keys', 'NSRA13', 'Pacific Islands', 'CONUS to NSRA', 'NonPeak', '$17.57', '$16.89');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US4965500', 'Florida Keys', 'NSRA13', 'Pacific Islands', 'CONUS to NSRA', 'Peak', '$20.73', '$19.93');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US68', 'Texas-South', 'NSRA2', 'Canada Central', 'CONUS to NSRA', 'NonPeak', '$27.63', '$44.24');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US68', 'Texas-South', 'NSRA2', 'Canada Central', 'CONUS to NSRA', 'Peak', '$32.60', '$52.20');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US68', 'Texas-South', 'NSRA13', 'Pacific Islands', 'CONUS to NSRA', 'NonPeak', '$10.65', '$16.89');
INSERT INTO public.stage_non_standard_locn_prices (origin_id, origin_area, destination_id, destination_area, move_type, season, hhg_price, ub_price) VALUES ('US68', 'Texas-South', 'NSRA13', 'Pacific Islands', 'CONUS to NSRA', 'Peak', '$12.57', '$19.93');


--
-- Data for Name: stage_oconus_to_conus_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'US16', 'Connecticut', 'NonPeak', '$29.28', '$33.75');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'US16', 'Connecticut', 'Peak', '$34.55', '$39.82');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'US47', 'Alabama', 'NonPeak', '$27.63', '$32.26');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'US47', 'Alabama', 'Peak', '$32.60', '$38.07');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'US4965500', 'Florida Keys', 'NonPeak', '$16.05', '$43.94');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'US4965500', 'Florida Keys', 'Peak', '$18.94', '$51.85');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'US68', 'Texas-South', 'NonPeak', '$17.57', '$44.51');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'US68', 'Texas-South', 'Peak', '$20.73', '$52.52');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'US16', 'Connecticut', 'NonPeak', '$15.28', '$31.91');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'US16', 'Connecticut', 'Peak', '$18.03', '$37.65');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'US47', 'Alabama', 'NonPeak', '$30.90', '$44.28');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'US47', 'Alabama', 'Peak', '$36.46', '$52.25');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'US4965500', 'Florida Keys', 'NonPeak', '$17.57', '$34.33');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'US4965500', 'Florida Keys', 'Peak', '$20.73', '$40.51');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'US68', 'Texas-South', 'NonPeak', '$24.38', '$34.90');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'US68', 'Texas-South', 'Peak', '$28.77', '$41.18');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'US4965500', 'Florida Keys', 'NonPeak', '$27.63', '$33.98');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'US4965500', 'Florida Keys', 'Peak', '$32.60', '$40.10');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'US16', 'Connecticut', 'NonPeak', '$24.38', '$44.62');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'US16', 'Connecticut', 'Peak', '$28.77', '$52.65');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'US47', 'Alabama', 'NonPeak', '$16.05', '$33.75');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'US47', 'Alabama', 'Peak', '$18.94', '$39.82');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'US68', 'Texas-South', 'NonPeak', '$17.57', '$34.45');
INSERT INTO public.stage_oconus_to_conus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_domestic_price_area_area, destination_domestic_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'US68', 'Texas-South', 'Peak', '$20.73', '$40.65');


--
-- Data for Name: stage_oconus_to_oconus_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$10.65', '$11.86');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'AS11', 'New South Wales/Australian Capital Territory', 'Peak', '$12.57', '$13.99');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'GE', 'Germany', 'NonPeak', '$8.88', '$12.99');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'GE', 'Germany', 'Peak', '$10.48', '$15.33');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'US8101000', 'Alaska (Zone) I', 'NonPeak', '$10.35', '$11.86');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('AS11', 'New South Wales/Australian Capital Territory', 'US8101000', 'Alaska (Zone) I', 'Peak', '$12.21', '$13.99');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$9.77', '$17.17');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'AS11', 'New South Wales/Australian Capital Territory', 'Peak', '$11.53', '$20.26');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'GE', 'Germany', 'NonPeak', '$11.35', '$16.89');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'GE', 'Germany', 'Peak', '$13.39', '$19.93');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'US8101000', 'Alaska (Zone) I', 'NonPeak', '$10.21', '$17.17');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('GE', 'Germany', 'US8101000', 'Alaska (Zone) I', 'Peak', '$12.05', '$20.26');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'AS11', 'New South Wales/Australian Capital Territory', 'NonPeak', '$11.35', '$11.66');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'AS11', 'New South Wales/Australian Capital Territory', 'Peak', '$13.39', '$13.76');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'GE', 'Germany', 'NonPeak', '$10.21', '$12.36');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'GE', 'Germany', 'Peak', '$12.05', '$14.58');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'US8101000', 'Alaska (Zone) I', 'NonPeak', '$11.35', '$17.59');
INSERT INTO public.stage_oconus_to_oconus_prices (origin_intl_price_area_id, origin_intl_price_area, destination_intl_price_area_id, destination_intl_price_area, season, hhg_shipping_linehaul_price, ub_price) VALUES ('US8101000', 'Alaska (Zone) I', 'US8101000', 'Alaska (Zone) I', 'Peak', '$13.39', '$19.93');


--
-- Data for Name: stage_other_intl_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('AS11', 'New South Wales/Australian Capital Territory', '$63.33', '$7.52', '$61.05', '$7.24', '$3.42', '$0.10', '$156.26', '$234.40', 'NonPeak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('AS11', 'New South Wales/Australian Capital Territory', '$74.10', '$9.15', '$71.43', '$8.47', '$4.00', '$0.12', '$182.82', '$274.25', 'Peak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('GE', 'Germany', '$57.76', '$7.52', '$86.31', '$7.24', '$4.85', '$0.16', '$172.00', '$258.00', 'NonPeak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('GE', 'Germany', '$67.58', '$9.15', '$100.98', '$8.47', '$5.67', '$0.19', '$201.24', '$301.86', 'Peak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('US16', 'Connecticut', '$67.48', '$7.52', '$63.88', '$7.24', '$5.39', '$0.18', '$156.26', '$217.96', 'NonPeak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('US16', 'Connecticut', '$78.95', '$9.15', '$74.74', '$8.47', '$6.31', '$0.21', '$182.82', '$255.01', 'Peak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('US47', 'Alabama', '$69.97', '$7.52', '$72.50', '$7.24', '$3.28', '$0.16', '$156.26', '$258.00', 'NonPeak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('US47', 'Alabama', '$81.86', '$9.15', '$84.82', '$8.47', '$3.84', '$0.19', '$182.82', '$301.86', 'Peak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('US4965500', 'Florida Keys', '$76.04', '$7.52', '$79.40', '$7.24', '$4.85', '$0.26', '$156.26', '$258.00', 'NonPeak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('US4965500', 'Florida Keys', '$88.97', '$9.15', '$92.90', '$8.47', '$5.67', '$0.30', '$182.82', '$301.86', 'Peak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('US68', 'Texas-South', '$69.97', '$7.52', '$72.50', '$7.24', '$4.33', '$0.12', '$145.31', '$258.00', 'NonPeak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('US68', 'Texas-South', '$81.86', '$9.15', '$84.82', '$8.47', '$5.07', '$0.14', '$170.01', '$301.86', 'Peak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('US8101000', 'Alaska (Zone) I', '$61.05', '$7.52', '$57.76', '$7.24', '$6.07', '$0.14', '$164.17', '$217.96', 'NonPeak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('US8101000', 'Alaska (Zone) I', '$71.43', '$9.15', '$67.58', '$8.47', '$7.10', '$0.16', '$192.08', '$255.01', 'Peak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('NSRA2', 'Canada Central', '$79.40', '$7.52', '$76.66', '$7.24', '$6.07', '$0.17', '$156.26', '$246.25', 'NonPeak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('NSRA2', 'Canada Central', '$92.90', '$9.15', '$89.69', '$8.47', '$7.10', '$0.20', '$182.82', '$288.11', 'Peak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('NSRA13', 'Pacific Islands', '$86.31', '$7.52', '$83.34', '$7.24', '$5.62', '$0.28', '$164.17', '$246.25', 'NonPeak');
INSERT INTO public.stage_other_intl_prices (rate_area_code, rate_area_name, hhg_origin_pack_price, hhg_destination_unpack_price, ub_origin_pack_price, ub_destination_unpack_price, origin_destination_sit_first_day_warehouse, origin_destination_sit_addl_days, sit_lte_50_miles, sit_gt_50_miles, season) VALUES ('NSRA13', 'Pacific Islands', '$100.98', '$9.15', '$97.51', '$8.47', '$6.58', '$0.33', '$192.08', '$288.11', 'Peak');


--
-- Data for Name: stage_price_escalation_discounts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_price_escalation_discounts (contract_year, forecasting_adjustment, discount, price_escalation) VALUES ('Base Period Year 1', '1.0000', '', '1');
INSERT INTO public.stage_price_escalation_discounts (contract_year, forecasting_adjustment, discount, price_escalation) VALUES ('Base Period Year 2', '1.0206', '', '1.0206');
INSERT INTO public.stage_price_escalation_discounts (contract_year, forecasting_adjustment, discount, price_escalation) VALUES ('Base Period Year 3', '1.0197', '', '1.0197');
INSERT INTO public.stage_price_escalation_discounts (contract_year, forecasting_adjustment, discount, price_escalation) VALUES ('Option Period 1', '1.0214', '', '1.0214');
INSERT INTO public.stage_price_escalation_discounts (contract_year, forecasting_adjustment, discount, price_escalation) VALUES ('Option Period 2', '1.0211', '', '1.0211');
INSERT INTO public.stage_price_escalation_discounts (contract_year, forecasting_adjustment, discount, price_escalation) VALUES ('Award Term 1', '1.0199', '', '1.0199');
INSERT INTO public.stage_price_escalation_discounts (contract_year, forecasting_adjustment, discount, price_escalation) VALUES ('Award Term 2', '1.0194', '', '1.0194');
INSERT INTO public.stage_price_escalation_discounts (contract_year, forecasting_adjustment, discount, price_escalation) VALUES ('Option Period 3', '1.0202', '', '1.0202');


--
-- Data for Name: stage_shipment_management_services_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_shipment_management_services_prices (contract_year, price_per_task_order) VALUES ('Base Period Year 1', '451.15');
INSERT INTO public.stage_shipment_management_services_prices (contract_year, price_per_task_order) VALUES ('Base Period Year 2', '454.23');
INSERT INTO public.stage_shipment_management_services_prices (contract_year, price_per_task_order) VALUES ('Base Period Year 3', '457.98');
INSERT INTO public.stage_shipment_management_services_prices (contract_year, price_per_task_order) VALUES ('Option Period 1', '459.85');
INSERT INTO public.stage_shipment_management_services_prices (contract_year, price_per_task_order) VALUES ('Option Period 2', '500.23');
INSERT INTO public.stage_shipment_management_services_prices (contract_year, price_per_task_order) VALUES ('Award Term 1', '500.69');
INSERT INTO public.stage_shipment_management_services_prices (contract_year, price_per_task_order) VALUES ('Award Term 2', '501.15');
INSERT INTO public.stage_shipment_management_services_prices (contract_year, price_per_task_order) VALUES ('Option Period 3', '502.19');


--
-- Data for Name: stage_transition_prices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.stage_transition_prices (contract_year, price_total_cost) VALUES ('One-time Amount (applied in Base Period Year 1)', '$ 450.18');


--
-- Name: TABLE stage_conus_to_oconus_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_conus_to_oconus_prices TO master;


--
-- Name: TABLE stage_counseling_services_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_counseling_services_prices TO master;


--
-- Name: TABLE stage_domestic_international_additional_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_domestic_international_additional_prices TO master;


--
-- Name: TABLE stage_domestic_linehaul_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_domestic_linehaul_prices TO master;


--
-- Name: TABLE stage_domestic_move_accessorial_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_domestic_move_accessorial_prices TO master;


--
-- Name: TABLE stage_domestic_service_area_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_domestic_service_area_prices TO master;


--
-- Name: TABLE stage_domestic_service_areas; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_domestic_service_areas TO master;


--
-- Name: TABLE stage_international_move_accessorial_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_international_move_accessorial_prices TO master;


--
-- Name: TABLE stage_international_service_areas; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_international_service_areas TO master;


--
-- Name: TABLE stage_non_standard_locn_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_non_standard_locn_prices TO master;


--
-- Name: TABLE stage_oconus_to_conus_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_oconus_to_conus_prices TO master;


--
-- Name: TABLE stage_oconus_to_oconus_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_oconus_to_oconus_prices TO master;


--
-- Name: TABLE stage_other_intl_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_other_intl_prices TO master;


--
-- Name: TABLE stage_price_escalation_discounts; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_price_escalation_discounts TO master;


--
-- Name: TABLE stage_shipment_management_services_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_shipment_management_services_prices TO master;


--
-- Name: TABLE stage_transition_prices; Type: ACL; Schema: public; Owner: -
--

GRANT ALL ON TABLE public.stage_transition_prices TO master;


--
-- PostgreSQL database dump complete
--

