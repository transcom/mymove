import React, { Component, createRef } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { push } from 'connected-react-router';

import Alert from 'shared/Alert';
import { withContext } from 'shared/AppContext';
import scrollToTop from 'shared/scrollToTop';
import { getOrdersForServiceMember, patchOrders, createUploadForDocument, deleteUpload } from 'services/internalApi';
import { updateOrders as updateOrdersAction } from 'store/entities/actions';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectMoveIsApproved,
  selectUploadsForCurrentOrders,
  selectHasCurrentPPM,
  selectEntitlementsForLoggedInUser,
} from 'store/entities/selectors';
import EditOrdersForm from 'components/Customer/EditOrdersForm/EditOrdersForm';
import { OrdersShape, HistoryShape } from 'types/customerShapes';
import { EntitlementShape, ExistingUploadsShape } from 'types';
import 'scenes/Review/Review.css';

class EditOrders extends Component {
  constructor(props) {
    super(props);

    this.filePondEl = createRef();
  }

  componentDidMount() {
    const { serviceMemberId, updateOrders } = this.props;
    getOrdersForServiceMember(serviceMemberId).then((response) => {
      updateOrders(response);
    });
  }

  handleUploadFile = (file) => {
    const { currentOrders } = this.props;
    const documentId = currentOrders?.uploaded_orders?.id;
    return createUploadForDocument(file, documentId);
  };

  handleUploadComplete = () => {
    const { serviceMemberId, updateOrders } = this.props;
    this.filePondEl.current?.removeFiles();
    return getOrdersForServiceMember(serviceMemberId).then((response) => {
      updateOrders(response);
    });
  };

  handleDeleteFile = (uploadId) => {
    const { serviceMemberId, updateOrders } = this.props;

    return deleteUpload(uploadId).then(() => {
      getOrdersForServiceMember(serviceMemberId).then((response) => {
        updateOrders(response);
      });
    });
  };

  submitOrders = (fieldValues) => {
    const { setFlashMessage, entitlement, updateOrders, currentOrders, spouseHasProGear, history } = this.props;

    let entitlementCouldChange = false;

    const fromFormSpouseHasProGear = (fieldValues.has_dependents && fieldValues.spouse_has_pro_gear) || false;

    if (fieldValues.has_dependents !== currentOrders.has_dependents || fromFormSpouseHasProGear !== spouseHasProGear) {
      entitlementCouldChange = true;
    }

    const newDutyStationId = fieldValues.new_duty_station.id;
    return patchOrders({
      ...fieldValues,
      new_duty_station_id: newDutyStationId,
      spouse_has_pro_gear: fromFormSpouseHasProGear,
    })
      .then((response) => {
        updateOrders(response);

        if (entitlementCouldChange) {
          setFlashMessage(
            'EDIT_ORDERS_SUCCESS',
            'info',
            `Your weight entitlement is now ${entitlement.sum.toLocaleString()} lbs.`,
            'Your changes have been saved. Note that the entitlement has also changed.',
          );
        } else {
          setFlashMessage('EDIT_ORDERS_SUCCESS', 'success', '', 'Your changes have been saved.');
        }

        history.goBack();
      })
      .catch(() => {
        scrollToTop();
      });
  };

  render() {
    const { error, schema, currentOrders, existingUploads, moveIsApproved } = this.props;
    return (
      <div className="usa-grid">
        {error && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {error.message}
            </Alert>
          </div>
        )}
        {moveIsApproved && (
          <div className="usa-width-one-whole error-message">
            <Alert type="warning" heading="Your move is approved">
              To make a change to your orders, you will need to contact your local PPPO office.
            </Alert>
          </div>
        )}
        {!moveIsApproved && (
          <div className="usa-width-one-whole">
            <EditOrdersForm
              initialValues={currentOrders}
              onSubmit={this.submitOrders}
              schema={schema}
              filePondEl={this.filePondEl}
              createUpload={this.handleUploadFile}
              onUploadComplete={this.handleUploadComplete}
              existingUploads={existingUploads}
              onDelete={this.handleDeleteFile}
            />
          </div>
        )}
      </div>
    );
  }
}

EditOrders.propTypes = {
  moveIsApproved: PropTypes.bool.isRequired,
  serviceMemberId: PropTypes.string.isRequired,
  setFlashMessage: PropTypes.func.isRequired,
  updateOrders: PropTypes.func.isRequired,
  currentOrders: OrdersShape.isRequired,
  history: HistoryShape.isRequired,
  entitlement: EntitlementShape.isRequired,
  existingUploads: ExistingUploadsShape,
  error: PropTypes.shape({
    message: PropTypes.string.isRequired,
  }),
  schema: PropTypes.shape({}),
  spouseHasProGear: PropTypes.bool,
};

EditOrders.defaultProps = {
  existingUploads: [],
  error: null,
  spouseHasProGear: false,
  schema: {},
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;
  const currentOrders = selectCurrentOrders(state) || {};
  const uploads = selectUploadsForCurrentOrders(state);

  return {
    currentOrders,
    serviceMemberId,
    existingUploads: uploads,
    error: get(state, 'orders.error'),
    moveIsApproved: selectMoveIsApproved(state),
    isPpm: selectHasCurrentPPM(state),
    schema: get(state, 'swaggerInternal.spec.definitions.CreateUpdateOrders', {}),
    entitlement: selectEntitlementsForLoggedInUser(state),
  };
}

const mapDispatchToProps = {
  push,
  updateOrders: updateOrdersAction,
  setFlashMessage: setFlashMessageAction,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(EditOrders));
