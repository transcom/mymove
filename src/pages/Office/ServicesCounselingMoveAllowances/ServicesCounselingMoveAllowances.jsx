/* eslint-disable camelcase */
import React, { useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import * as Yup from 'yup';

import AllowancesDetailForm from '../../../components/Office/AllowancesDetailForm/AllowancesDetailForm';

import styles from 'styles/documentViewerWithSidebar.module.scss';
import { milmoveLogger } from 'utils/milmoveLog';
import { ORDERS_BRANCH_OPTIONS, ORDERS_PAY_GRADE_TYPE, ORDERS_TYPE } from 'constants/orders';
import { ORDERS } from 'constants/queryKeys';
import { servicesCounselingRoutes } from 'constants/routes';
import { MOVE_STATUSES, FEATURE_FLAG_KEYS } from 'shared/constants';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { counselingUpdateAllowance } from 'services/ghcApi';
import { dropdownInputOptions } from 'utils/formatters';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const branchDropdownOption = dropdownInputOptions(ORDERS_BRANCH_OPTIONS);

const validationSchema = Yup.object({
  proGearWeight: Yup.number()
    .min(0, 'Pro-gear weight must be greater than or equal to 0')
    .max(2000, "Enter a weight that does not go over the customer's maximum allowance")
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .notRequired(),
  proGearWeightSpouse: Yup.number()
    .min(0, 'Spouse pro-gear weight must be greater than or equal to 0')
    .max(500, "Enter a weight that does not go over the customer's maximum allowance")
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .notRequired(),
  gunSafeWeight: Yup.number()
    .min(0, 'Gun safe weight must be greater than or equal to 0')
    .max(500, "Enter a weight that does not go over the customer's maximum allowance")
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .notRequired(),
  requiredMedicalEquipmentWeight: Yup.number()
    .min(0, 'Required medical equipment weight must be greater than or equal to 0')
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .notRequired(),
  storageInTransit: Yup.number()
    .min(0, 'Storage in transit (days) must be greater than or equal to 0')
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .notRequired(),
  weightRestriction: Yup.number()
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .when('adminRestrictedWeightLocation', {
      is: true,
      then: (schema) =>
        schema
          .min(1, 'Weight restriction must be greater than 0')
          .max(18000, 'Weight restriction cannot exceed 18,000 lbs')
          .required('Weight restriction is required when Admin Restricted Weight Location is enabled'),
      otherwise: (schema) => schema.notRequired().nullable(),
    }),
  adminRestrictedWeightLocation: Yup.boolean().notRequired(),
  ubWeightRestriction: Yup.number()
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .when('adminRestrictedUBWeightLocation', {
      is: true,
      then: (schema) =>
        schema
          .min(1, 'UB weight restriction must be greater than 0')
          .max(2000, 'UB weight restriction cannot exceed 2,000 lbs')
          .required('UB weight restriction is required when Admin Restricted UB Weight Location is enabled'),
      otherwise: (schema) => schema.notRequired().nullable(),
    }),
  adminRestrictedUBWeightLocation: Yup.boolean().notRequired(),
  ubAllowance: Yup.number()
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .min(0, 'UB weight allowance must be 0 or more')
    .max(2000, 'UB weight allowance cannot exceed 2,000 lbs.'),
});
const ServicesCounselingMoveAllowances = () => {
  const { moveCode } = useParams();
  const navigate = useNavigate();

  const { move, orders, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const [isGunSafeEnabled, setIsGunSafeEnabled] = useState(false);
  const orderId = move?.ordersId;

  useEffect(() => {
    const fetchData = async () => {
      setIsGunSafeEnabled(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.GUN_SAFE));
    };
    fetchData();
  }, []);

  const handleClose = () => {
    navigate(`../${servicesCounselingRoutes.MOVE_VIEW_PATH}`);
  };
  const queryClient = useQueryClient();
  const { mutate: mutateOrders } = useMutation(counselingUpdateAllowance, {
    onSuccess: (data, variables) => {
      const updatedOrder = data.orders[variables.orderID];
      queryClient.setQueryData([ORDERS, variables.orderID], {
        orders: {
          [`${variables.orderID}`]: updatedOrder,
        },
      });
      queryClient.invalidateQueries(ORDERS);
      handleClose();
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const order = Object.values(orders)?.[0];
  const onSubmit = (values) => {
    const {
      grade,
      agency,
      proGearWeight,
      proGearWeightSpouse,
      requiredMedicalEquipmentWeight,
      organizationalClothingAndIndividualEquipment,
      storageInTransit,
      gunSafe,
      gunSafeWeight,
      adminRestrictedWeightLocation,
      weightRestriction,
      adminRestrictedUBWeightLocation,
      ubWeightRestriction,
      accompaniedTour,
      dependentsTwelveAndOver,
      dependentsUnderTwelve,
    } = values;
    const body = {
      issueDate: order.date_issued,
      newDutyLocationId: order.destinationDutyLocation.id,
      ordersNumber: order.order_number,
      ordersType: order.order_type,
      originDutyLocationId: order.originDutyLocation.id,
      reportByDate: order.report_by_date,
      grade,
      agency,
      proGearWeight: Number(proGearWeight),
      proGearWeightSpouse: Number(proGearWeightSpouse),
      requiredMedicalEquipmentWeight: Number(requiredMedicalEquipmentWeight),
      storageInTransit: Number(storageInTransit),
      organizationalClothingAndIndividualEquipment,
      gunSafe,
      weightRestriction: adminRestrictedWeightLocation && weightRestriction ? Number(weightRestriction) : null,
      ubWeightRestriction: adminRestrictedUBWeightLocation && ubWeightRestriction ? Number(ubWeightRestriction) : null,
      accompaniedTour,
      dependentsTwelveAndOver: Number(dependentsTwelveAndOver),
      dependentsUnderTwelve: Number(dependentsUnderTwelve),
      ubAllowance: Number(values.ubAllowance),
    };
    if (isGunSafeEnabled) body.gunSafeWeight = Number(gunSafeWeight);

    return mutateOrders({ orderID: orderId, ifMatchETag: order.eTag, body });
  };

  const counselorCanEdit =
    move.status === MOVE_STATUSES.NEEDS_SERVICE_COUNSELING ||
    move.status === MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED ||
    (move.status === MOVE_STATUSES.APPROVALS_REQUESTED && !move.availableToPrimeAt); // status is set to 'Approval Requested' if customer uploads amended orders.

  const { entitlement, grade, agency } = order;
  const {
    proGearWeight,
    proGearWeightSpouse,
    requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment,
    gunSafe,
    gunSafeWeight,
    weightRestriction,
    ubWeightRestriction,
    storageInTransit,
    dependentsUnderTwelve,
    dependentsTwelveAndOver,
    accompaniedTour,
  } = entitlement;

  const initialValues = {
    grade,
    agency,
    proGearWeight: `${proGearWeight}`,
    proGearWeightSpouse: `${proGearWeightSpouse}`,
    requiredMedicalEquipmentWeight: `${requiredMedicalEquipmentWeight}`,
    storageInTransit: `${storageInTransit}`,
    gunSafe,
    gunSafeWeight: `${gunSafeWeight}`,
    adminRestrictedWeightLocation: weightRestriction > 0,
    weightRestriction: weightRestriction ? `${weightRestriction}` : '0',
    adminRestrictedUBWeightLocation: ubWeightRestriction > 0,
    ubWeightRestriction: ubWeightRestriction ? `${ubWeightRestriction}` : '0',
    organizationalClothingAndIndividualEquipment,
    accompaniedTour,
    dependentsUnderTwelve: `${dependentsUnderTwelve}`,
    dependentsTwelveAndOver: `${dependentsTwelveAndOver}`,
  };

  const civilianTDYUBMove =
    order.order_type === ORDERS_TYPE.TEMPORARY_DUTY &&
    order.grade === ORDERS_PAY_GRADE_TYPE.CIVILIAN_EMPLOYEE &&
    (order.originDutyLocation?.address?.isOconus || order.destinationDutyLocation?.address?.isOconus);

  return (
    <div className={styles.sidebar}>
      <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
        {(formik) => (
          <form onSubmit={formik.handleSubmit}>
            <div className={styles.content}>
              <div className={styles.top}>
                <Button
                  className={styles.closeButton}
                  data-testid="closeSidebar"
                  type="button"
                  onClick={handleClose}
                  unstyled
                >
                  <FontAwesomeIcon icon="times" title="Close sidebar" aria-label="Close sidebar" />
                </Button>
                <h2 className={styles.header} data-testid="allowances-header">
                  View allowances
                </h2>
                <div>
                  <Link className={styles.viewAllowances} data-testid="view-orders" to="../orders">
                    View orders
                  </Link>
                </div>
              </div>
              <div className={styles.body}>
                <AllowancesDetailForm
                  entitlements={order.entitlement}
                  branchOptions={branchDropdownOption}
                  header="Counseling"
                  civilianTDYUBMove={civilianTDYUBMove}
                  formIsDisabled={!counselorCanEdit}
                />
              </div>
              <div className={styles.bottom}>
                <div className={styles.buttonGroup}>
                  <Button
                    disabled={formik.isSubmitting || !formik.isValid || !counselorCanEdit}
                    data-testid="scAllowancesSave"
                    type="submit"
                  >
                    Save
                  </Button>
                  <Button type="button" secondary onClick={handleClose}>
                    Cancel
                  </Button>
                </div>
              </div>
            </div>
          </form>
        )}
      </Formik>
    </div>
  );
};

export default ServicesCounselingMoveAllowances;
