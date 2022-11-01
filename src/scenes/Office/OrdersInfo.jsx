import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert';
import DocumentContent from 'shared/DocumentViewer/DocumentContent';
import OrdersViewerPanel from './OrdersViewerPanel';
import { getRequestStatus } from 'shared/Swagger/selectors';
import { loadMove, selectMove } from 'shared/Entities/modules/moves';
import { loadOrders, loadOrdersLabel } from 'shared/Entities/modules/orders';
import { loadServiceMember, selectServiceMember } from 'shared/Entities/modules/serviceMembers';
import { stringifyName } from 'shared/utils/serviceMember';
import { selectOrdersById } from 'store/entities/selectors';

import './office.scss';

class OrdersInfo extends Component {
  componentDidMount() {
    this.props.loadMove(this.props.moveId);
  }

  componentDidUpdate(prevProps) {
    const { serviceMemberId, ordersId } = this.props;
    // Both serviceMemberId and ordersId come from the move
    // If we have one, we can safely assume we have the other
    if (serviceMemberId !== prevProps.serviceMemberId) {
      this.props.loadServiceMember(serviceMemberId);
      this.props.loadOrders(ordersId);
    }
  }

  render() {
    const { serviceMember, uploads } = this.props;
    const name = stringifyName(serviceMember);
    document.title = `Orders for ${name}`;

    if (!this.props.loadDependenciesHasSuccess && !this.props.loadDependenciesHasError) return <LoadingPlaceholder />;
    if (this.props.loadDependenciesHasError)
      return (
        <div className="grid-container-widescreen usa-prose">
          <div className="grid-row">
            <div className="grid-col-12 error-message">
              <Alert type="error" heading="An error occurred">
                Something went wrong contacting the server.
              </Alert>
            </div>
          </div>
        </div>
      );
    return (
      <div className="grid-container-widescreen usa-prose">
        <div className="grid-row grid-gap">
          <div className="grid-col-8 document-contents">
            {uploads.map((upload) => (
              <DocumentContent
                key={upload.url}
                url={upload.url}
                filename={upload.filename}
                contentType={upload.contentType}
                status={upload.status}
              />
            ))}
          </div>
          <div className="grid-col-4 orders-page-fields">
            <OrdersViewerPanel title={name} className="document-viewer" moveId={this.props.match.params.moveId} />
          </div>
        </div>
      </div>
    );
  }
}

OrdersInfo.propTypes = {
  loadMove: PropTypes.func.isRequired,
};

const mapStateToProps = (state, ownProps) => {
  const moveId = ownProps.match.params.moveId;
  const move = selectMove(state, moveId);
  const ordersId = move.orders_id;
  const orders = selectOrdersById(state, ordersId);
  const uploads = orders?.uploaded_orders?.uploads || [];
  const serviceMemberId = move.service_member_id;
  const serviceMember = selectServiceMember(state, serviceMemberId);
  const loadOrdersRequest = getRequestStatus(state, loadOrdersLabel);

  return {
    swaggerError: state.swaggerInternal.hasErrored,
    moveId,
    ordersSchema: get(state, 'swaggerInternal.spec.definitions.CreateUpdateOrders', {}),
    ordersId,
    serviceMember,
    serviceMemberId,
    uploads,
    loadDependenciesHasSuccess: loadOrdersRequest.isSuccess,
    loadDependenciesHasError: loadOrdersRequest.error,
  };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({ loadMove, loadOrders, loadServiceMember }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(OrdersInfo);
