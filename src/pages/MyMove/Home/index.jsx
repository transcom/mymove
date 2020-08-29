/* eslint-disable no-console */
import React, { Component } from 'react';
import { isNil, isEmpty } from 'lodash';
import { func, arrayOf, shape, string, objectOf, object } from 'prop-types';
import { connect } from 'react-redux';
import { Alert } from '@trussworks/react-uswds';

import styles from './Home.module.scss';

import Helper from 'components/Customer/Home/Helper';
import Step from 'components/Customer/Home/Step';
import DocsUploaded from 'components/Customer/Home/DocsUploaded';
import ShipmentList from 'components/Customer/Home/ShipmentList';
import Contact from 'components/Customer/Home/Contact';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { selectActiveOrLatestOrdersFromEntities, selectUploadedOrders } from 'shared/Entities/modules/orders';
import { selectActiveOrLatestMove } from 'shared/Entities/modules/moves';

const shipments = [
  { type: SHIPMENT_OPTIONS.PPM, id: '#123ABC-001' },
  { type: SHIPMENT_OPTIONS.HHG, id: '#123ABC-002' },
  { type: SHIPMENT_OPTIONS.NTS, id: '#123ABC-003' },
];

class Home extends Component {
  componentDidMount() {
    const { showLoggedInUser } = this.props;
    showLoggedInUser();
  }

  handleShipmentClick = (shipment) => {
    console.log('this is the shipment', shipment);
  };

  handleNewPathClick = (path) => {
    const { push } = this.props;
    push(path);
  };

  checkOrdersCompleted = () => {
    const { orders, uploadedOrderDocuments } = this.props;
    return (orders?.['uploaded_orders']?.uploads || []).length > 0 || uploadedOrderDocuments.length > 0;
  };

  render() {
    const { move } = this.props;
    const ordersPath = '/orders/';
    const shipmentSelectionPath = `/moves/${move.id}/select-type`;
    const confirmationPath = `/moves/${move.id}/review`;
    const profileEditPath = '/moves/review/edit-profile';
    const ordersEditPath = `/moves/${move.id}/review/edit-orders`;
    const ordersCompleted = this.checkOrdersCompleted();
    const hasShipment = !isNil(move.personally_procured_moves) || !isNil(move.shipments);
    const hasSubmittedMove = !isEmpty(move) && move.status !== 'DRAFT';

    return (
      <div className={`usa-prose grid-container ${styles['grid-container']}`}>
        <header className={styles['customer-header']}>
          <h2>Riley Baker</h2>
          <p>
            You&apos;re leaving <strong>Buckley AFB</strong>
          </p>
        </header>
        <Alert className="margin-top-2 margin-bottom-2" slim type="success">
          Thank you for adding your Profile information
        </Alert>

        <Helper title="Next step: Add your orders">
          <ul>
            {[
              'If you have a hard copy, you can take photos of each page',
              'If you have a PDF, you can upload that',
            ].map((helpText) => (
              <li key={helpText}>
                <span>{helpText}</span>
              </li>
            ))}
          </ul>
        </Helper>
        <Step
          complete
          completedHeaderText="Profile complete"
          editBtnLabel="Edit"
          headerText="Profile complete"
          step="1"
          onEditBtnClick={() => this.handleNewPathClick(profileEditPath)}
        >
          <p className={styles.description}>Make sure to keep your personal information up to date during your move</p>
        </Step>

        <Step
          complete={ordersCompleted}
          completedHeaderText="Orders uploaded"
          editBtnLabel={ordersCompleted ? 'Edit' : ''}
          onEditBtnClick={() => this.handleNewPathClick(ordersEditPath)}
          headerText="Upload orders"
          actionBtnLabel={!ordersCompleted ? 'Add orders' : ''}
          onActionBtnClick={() => this.handleNewPathClick(ordersPath)}
          step="2"
        >
          {ordersCompleted && (
            <DocsUploaded
              files={[
                { filename: 'Screen Shot 2020-09-11 at 12.56.58 PM.png' },
                { filename: 'Screen Shot 2020-09-11 at 12.58.12 PM.png' },
                { filename: 'orderspage3_20200723.png' },
              ]}
            />
          )}
        </Step>

        <Step
          actionBtnLabel={hasShipment ? 'Add another shipment' : 'Plan your shipments'}
          actionBtnDisabled={!ordersCompleted}
          onActionBtnClick={() => this.handleNewPathClick(shipmentSelectionPath)}
          complete={hasShipment}
          completedHeaderText="Shipments"
          headerText="Shipments"
          secondaryBtn
          secondaryClassName="margin-top-2"
          step="3"
        >
          {hasShipment && <ShipmentList shipments={shipments} onShipmentClick={this.handleShipmentClick} />}
        </Step>

        <Step
          complete={hasSubmittedMove}
          actionBtnDisabled={!hasShipment}
          actionBtnLabel={!hasSubmittedMove ? 'Review and submit' : ''}
          containerClassName="margin-bottom-8"
          headerText="Confirm move request"
          completedHeaderText="Move request confirmed"
          onActionBtnClick={() => this.handleNewPathClick(confirmationPath)}
          step="4"
        >
          <p className={styles.description}>
            {hasSubmittedMove
              ? 'Move submitted.'
              : 'Review your move details and sign the legal paperwork, then send the info on to your move counselor.'}
          </p>
        </Step>
        <Contact
          header="Contacts"
          dutyStationName="Seymour Johnson AFB"
          officeType="Origin Transportation Office"
          telephone="(919) 722-5458"
        />
      </div>
    );
  }
}

Home.propTypes = {
  showLoggedInUser: func.isRequired,
  // eslint-disable-next-line react/no-unused-prop-types
  uploadedOrderDocuments: arrayOf(
    shape({
      filename: string.isRequired,
    }),
  ).isRequired,
  // eslint-disable-next-line react/forbid-prop-types
  orders: object.isRequired,
  // eslint-disable-next-line react/forbid-prop-types
  push: func.isRequired,
  move: objectOf(
    shape({
      id: string,
    }),
  ).isRequired,
};

const mapStateToProps = (state) => {
  return {
    uploadedOrderDocuments: selectUploadedOrders(state),
    // TODO: change when we support multiple moves
    move: selectActiveOrLatestMove(state),
    orders: selectActiveOrLatestOrdersFromEntities(state),
  };
};

// in order to avoid setting up proxy server only for storybook, pass in stub function so API requests don't fail
const mergeProps = (stateProps, dispatchProps, ownProps) => ({
  ...stateProps,
  ...dispatchProps,
  ...ownProps,
});

const mapDispatchToProps = {
  showLoggedInUser: showLoggedInUserAction,
};

export default connect(mapStateToProps, mapDispatchToProps, mergeProps)(Home);
