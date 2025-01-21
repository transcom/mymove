import React, { useEffect, useState } from 'react';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { generatePath, useLocation, useNavigate, useParams } from 'react-router-dom';
import { connect } from 'react-redux';

import styles from './ServicesCounselingAddOrders.module.scss';

import { constructSCOrderOconusFields, dropdownInputOptions, formatYesNoAPIValue } from 'utils/formatters';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';
import AddOrdersForm from 'components/Office/AddOrdersForm/AddOrdersForm';
import { counselingCreateOrder } from 'services/ghcApi';
import { ORDERS } from 'constants/queryKeys';
import { formatDateForSwagger } from 'shared/dates';
import { servicesCounselingRoutes } from 'constants/routes';
import { milmoveLogger } from 'utils/milmoveLog';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { elevatedPrivilegeTypes } from 'constants/userPrivileges';
import withRouter from 'utils/routing';
import { withContext } from 'shared/AppContext';
import { setCanAddOrders as setCanAddOrdersAction } from 'store/general/actions';
import { selectCanAddOrders } from 'store/entities/selectors';

const ServicesCounselingAddOrders = ({ userPrivileges, canAddOrders, setCanAddOrders }) => {
  const { customerId } = useParams();
  const { state } = useLocation();
  const isSafetyMoveSelected = state?.isSafetyMoveSelected;
  const isBluebarkMoveSelected = state?.isBluebarkMoveSelected;
  const navigate = useNavigate();
  const [isSafetyMoveFF, setSafetyMoveFF] = useState(false);
  const [hasSubmitted, setHasSubmitted] = useState(false);

  const handleBack = () => {
    navigate(-1);
  };

  const handleClose = (moveCode) => {
    const path = generatePath(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, {
      moveCode,
    });
    navigate(path);
  };

  const queryClient = useQueryClient();
  const { mutate: mutateOrders } = useMutation(counselingCreateOrder, {
    onSuccess: (data) => {
      setCanAddOrders(false);
      const orderID = Object.keys(data.orders)[0];
      const updatedOrder = data.orders[orderID];
      queryClient.setQueryData([ORDERS, orderID], {
        orders: {
          [`${orderID}`]: updatedOrder,
        },
      });

      queryClient.invalidateQueries(ORDERS);
      handleClose(updatedOrder.moveCode);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  useEffect(() => {
    const redirectUser = async () => {
      if (!canAddOrders && !hasSubmitted) {
        const path = generatePath(servicesCounselingRoutes.BASE_QUEUE_VIEW_PATH);
        navigate(path);
      }
    };
    redirectUser();
  }, [canAddOrders, hasSubmitted, navigate]);

  useEffect(() => {
    isBooleanFlagEnabled('safety_move').then((enabled) => {
      setSafetyMoveFF(enabled);
    });
  }, []);

  const isSafetyPrivileged =
    isSafetyMoveFF && isSafetyMoveSelected !== false
      ? userPrivileges?.some((privilege) => privilege.privilegeType === elevatedPrivilegeTypes.SAFETY)
      : false;

  const allowedOrdersTypes = {
    ...ORDERS_TYPE_OPTIONS,
    ...(isSafetyPrivileged ? { SAFETY: 'Safety' } : {}),
    ...(isBluebarkMoveSelected ? { BLUEBARK: 'BLUEBARK' } : {}),
  };
  const ordersTypeOptions = dropdownInputOptions(allowedOrdersTypes);

  const getInitialOrdersType = () => {
    if (isSafetyMoveSelected) {
      return 'SAFETY';
    }
    if (isBluebarkMoveSelected) {
      return 'BLUEBARK';
    }
    return '';
  };

  const initialValues = {
    ordersType: getInitialOrdersType(),
    issueDate: '',
    reportByDate: '',
    hasDependents: '',
    newDutyLocation: '',
    grade: '',
    originDutyLocation: '',
    accompaniedTour: '',
    dependentsUnderTwelve: '',
    dependentsTwelveAndOver: '',
  };

  const handleSubmit = (values) => {
    setHasSubmitted(true);
    const oconusFields = constructSCOrderOconusFields(values);
    const body = {
      ...values,
      serviceMemberId: customerId,
      newDutyLocationId: values.newDutyLocation.id,
      hasDependents: formatYesNoAPIValue(values.hasDependents),
      reportByDate: formatDateForSwagger(values.reportByDate),
      issueDate: formatDateForSwagger(values.issueDate),
      grade: values.grade,
      originDutyLocationId: values.originDutyLocation.id,
      spouseHasProGear: false,
      ...oconusFields,
    };
    mutateOrders({ body });
  };

  return (
    <GridContainer data-testid="main-container">
      <Grid row className={styles.ordersFormContainer} data-testid="orders-form-container">
        <Grid col>
          <AddOrdersForm
            onSubmit={handleSubmit}
            ordersTypeOptions={ordersTypeOptions}
            initialValues={initialValues}
            onBack={handleBack}
            isSafetyMoveSelected={isSafetyMoveSelected}
            isBluebarkMoveSelected={isBluebarkMoveSelected}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};
const mapStateToProps = (state) => {
  const canAddOrders = selectCanAddOrders(state);

  return { canAddOrders };
};

const mapDispatchToProps = { setCanAddOrders: setCanAddOrdersAction };

export default withContext(withRouter(connect(mapStateToProps, mapDispatchToProps)(ServicesCounselingAddOrders)));
