import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { useNavigate, useParams } from 'react-router';

import { ORDERS_TYPE_OPTIONS } from 'constants/orders';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { patchOrders, getResponseError, getOrders } from 'services/internalApi';
import { updateOrders as updateOrdersAction } from 'store/entities/actions';
import { withContext } from 'shared/AppContext';
import { formatDateForSwagger } from 'shared/dates';
import { formatYesNoInputValue, formatYesNoAPIValue, dropdownInputOptions } from 'utils/formatters';
import { selectServiceMemberFromLoggedInUser, selectOrdersForLoggedInUser } from 'store/entities/selectors';
import { generalRoutes } from 'constants/routes';
import OrdersInfoForm from 'components/Customer/OrdersInfoForm/OrdersInfoForm';
import withRouter from 'utils/routing';

const Orders = ({ context, serviceMemberId, updateOrders, orders }) => {
  const [serverError, setServerError] = useState(null);
  const navigate = useNavigate();
  const { orderId } = useParams();
  const currentOrders = orders.find((order) => order.id === orderId);

  const handleBack = () => {
    navigate(generalRoutes.HOME_PATH);
  };

  const handleNext = (id) => {
    navigate(`/orders/upload/${id}`);
  };

  const submitOrders = async (values) => {
    const { rank, ...pendingValues } = {
      ...values,
      service_member_id: serviceMemberId,
      new_duty_location_id: values.new_duty_location.id,
      has_dependents: formatYesNoAPIValue(values.has_dependents),
      report_by_date: formatDateForSwagger(values.report_by_date),
      issue_date: formatDateForSwagger(values.issue_date),
      grade: values.grade,
      rankShortName: values.rank,
      origin_duty_location_id: values.origin_duty_location.id,
      spouse_has_pro_gear: false,
    };

    if (currentOrders?.id) {
      try {
        pendingValues.id = currentOrders.id;
        await patchOrders(pendingValues);
        await getOrders(currentOrders.id).then((response) => {
          updateOrders(response);
        });
        handleNext(currentOrders.id);
      } catch (error) {
        const { response } = error;
        const errorMessage = getResponseError(response, 'failed to update/create orders due to server error');
        setServerError(errorMessage);
      }
    }
  };

  const initialValues = {
    orders_type: currentOrders?.orders_type || '',
    issue_date: currentOrders?.issue_date || '',
    report_by_date: currentOrders?.report_by_date || '',
    has_dependents: formatYesNoInputValue(currentOrders?.has_dependents),
    new_duty_location: currentOrders?.new_duty_location || null,
    grade: currentOrders?.grade || null,
    origin_duty_location: currentOrders?.origin_duty_location || null,
    rank: currentOrders?.payGradeRank?.rankShortName || '',
  };
  // Only allow PCS unless feature flag is on
  const showAllOrdersTypes = context.flags?.allOrdersTypes;
  const allowedOrdersTypes = showAllOrdersTypes
    ? ORDERS_TYPE_OPTIONS
    : { PERMANENT_CHANGE_OF_STATION: ORDERS_TYPE_OPTIONS.PERMANENT_CHANGE_OF_STATION };

  const ordersTypeOptions = dropdownInputOptions(allowedOrdersTypes);

  return (
    <GridContainer data-testid="main-container">
      <NotificationScrollToTop dependency={serverError} />

      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid row data-testid="orders-form-container">
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <OrdersInfoForm
            ordersTypeOptions={ordersTypeOptions}
            initialValues={initialValues}
            onSubmit={submitOrders}
            onBack={handleBack}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

Orders.propTypes = {
  context: PropTypes.shape({
    flags: PropTypes.shape({
      allOrdersTypes: PropTypes.bool,
    }).isRequired,
  }).isRequired,
  serviceMemberId: PropTypes.string.isRequired,
  updateOrders: PropTypes.func.isRequired,
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const orders = selectOrdersForLoggedInUser(state);

  return {
    serviceMemberId: serviceMember?.id,
    orders,
  };
};

const mapDispatchToProps = {
  updateOrders: updateOrdersAction,
};

export default withContext(withRouter(connect(mapStateToProps, mapDispatchToProps)(Orders)));
