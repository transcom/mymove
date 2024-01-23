import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import OrdersInfoForm from 'components/Customer/OrdersInfoForm/OrdersInfoForm';
import {
  getServiceMember,
  getOrdersForServiceMember,
  createOrders,
  patchOrders,
  patchServiceMember,
  getResponseError,
} from 'services/internalApi';
import {
  updateOrders as updateOrdersAction,
  updateServiceMember as updateServiceMemberAction,
} from 'store/entities/actions';
import { withContext } from 'shared/AppContext';
import { formatDateForSwagger } from 'shared/dates';
import { OrdersShape } from 'types/customerShapes';
import { formatYesNoInputValue, formatYesNoAPIValue, dropdownInputOptions } from 'utils/formatters';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { selectServiceMemberFromLoggedInUser, selectCurrentOrders } from 'store/entities/selectors';
import { DutyLocationShape, RouterShape } from 'types';
import { customerRoutes, generalRoutes } from 'constants/routes';
import withRouter from 'utils/routing';

export class Orders extends Component {
  constructor(props) {
    super(props);

    this.state = {
      isLoading: true,
      serverError: null,
    };
  }

  componentDidMount() {
    const { serviceMemberId, currentOrders, updateOrders } = this.props;

    if (!currentOrders) {
      this.setState({ isLoading: false });
    } else {
      getOrdersForServiceMember(serviceMemberId).then((response) => {
        updateOrders(response);
        this.setState({ isLoading: false });
      });
    }
  }

  render() {
    const {
      context,
      currentDutyLocation,
      router: { navigate },
      serviceMemberId,
      currentOrders,
      updateOrders,
      updateServiceMember,
    } = this.props;

    const { isLoading, serverError } = this.state;

    if (isLoading) return <LoadingPlaceholder />;

    const handleBack = () => {
      navigate(generalRoutes.HOME_PATH);
    };

    const handleNext = () => {
      navigate(customerRoutes.ORDERS_UPLOAD_PATH);
    };

    const submitOrders = (values) => {
      const pendingValues = {
        ...values,
        service_member_id: serviceMemberId,
        new_duty_location_id: values.new_duty_location.id,
        has_dependents: formatYesNoAPIValue(values.has_dependents),
        report_by_date: formatDateForSwagger(values.report_by_date),
        issue_date: formatDateForSwagger(values.issue_date),
        grade: values.grade,
        origin_duty_location_id: values.origin_duty_location.id,
        spouse_has_pro_gear: false,
      };

      const payload = {
        id: serviceMemberId,
        rank: values.grade,
        current_location_id: values.origin_duty_location.id,
      };

      patchServiceMember(payload)
        .then(updateServiceMember)
        .catch((e) => {
          // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
          const { response } = e;
          const errorMessage = getResponseError(response, 'failed to update service member due to server error');
          this.setState({ serverError: errorMessage });
        });

      if (currentOrders?.id) {
        pendingValues.id = currentOrders.id;
        return patchOrders(pendingValues)
          .then(updateOrders)
          .then(handleNext)
          .catch((e) => {
            // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
            const { response } = e;
            const errorMessage = getResponseError(response, 'failed to update orders due to server error');
            this.setState({ serverError: errorMessage });
          });
      }

      return createOrders(pendingValues)
        .then(updateOrders)
        .then(() => getServiceMember(serviceMemberId))
        .then(updateServiceMember)
        .then(handleNext)
        .catch((e) => {
          // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
          const { response } = e;
          const errorMessage = getResponseError(response, 'failed to create orders due to server error');
          this.setState({ serverError: errorMessage });
        });
    };

    const initialValues = {
      orders_type: currentOrders?.orders_type || '',
      issue_date: currentOrders?.issue_date || '',
      report_by_date: currentOrders?.report_by_date || '',
      has_dependents: formatYesNoInputValue(currentOrders?.has_dependents),
      new_duty_location: currentOrders?.new_duty_location || null,
      grade: currentOrders?.grade || null,
      origin_duty_location: currentOrders?.origin_duty_location || null,
    };

    // Only allow PCS unless feature flag is on
    const showAllOrdersTypes = context.flags?.allOrdersTypes;
    const allowedOrdersTypes = showAllOrdersTypes
      ? ORDERS_TYPE_OPTIONS
      : { PERMANENT_CHANGE_OF_STATION: ORDERS_TYPE_OPTIONS.PERMANENT_CHANGE_OF_STATION };

    const ordersTypeOptions = dropdownInputOptions(allowedOrdersTypes);

    return (
      <GridContainer>
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

        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <OrdersInfoForm
              ordersTypeOptions={ordersTypeOptions}
              initialValues={initialValues}
              currentDutyLocation={currentDutyLocation}
              onSubmit={submitOrders}
              onBack={handleBack}
            />
          </Grid>
        </Grid>
      </GridContainer>
    );
  }
}

Orders.propTypes = {
  context: PropTypes.shape({
    flags: PropTypes.shape({
      allOrdersTypes: PropTypes.bool,
    }).isRequired,
  }).isRequired,
  serviceMemberId: PropTypes.string.isRequired,
  currentOrders: OrdersShape,
  updateOrders: PropTypes.func.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentDutyLocation: DutyLocationShape,
  router: RouterShape,
};

Orders.defaultProps = {
  currentOrders: null,
  currentDutyLocation: {},
  router: {},
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const orders = selectCurrentOrders(state);

  return {
    serviceMemberId: serviceMember?.id,
    currentOrders: selectCurrentOrders(state),
    currentDutyLocation: orders?.origin_duty_location || {},
  };
};

const mapDispatchToProps = {
  updateOrders: updateOrdersAction,
  updateServiceMember: updateServiceMemberAction,
};

export default withContext(withRouter(connect(mapStateToProps, mapDispatchToProps)(Orders)));
