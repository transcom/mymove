/* eslint-disable no-console */
import React, { Component } from 'react';
import { func, arrayOf, shape, string, objectOf } from 'prop-types';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import { Alert } from '@trussworks/react-uswds';

import styles from './Home.module.scss';

import Helper from 'components/Customer/Home/Helper';
import Step from 'components/Customer/Home/Step';
import DocsUploaded from 'components/Customer/Home/DocsUploaded';
import ShipmentList from 'components/Customer/Home/ShipmentList';
import Contact from 'components/Customer/Home/Contact';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { selectUploadedOrders } from 'shared/Entities/modules/orders';
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

  handleCLickToOrders = () => {
    const { history } = this.props;
    const path = '/orders/';
    history.push(path);
  };

  handleCLickToSelectShipment = () => {
    const { history, move } = this.props;
    const path = `/moves/${move.id}/select-type`;
    history.push(path);
  };

  render() {
    const ordersCompleted = false;
    const hasShipment = false;
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
          editBtnDisabled
          editBtnLabel="Edit"
          headerText="Profile complete"
          step="1"
          onEditClick={(e) => {
            e.preventDefault();
            console.log('edit clicked');
          }}
        >
          <p className={styles.description}>Make sure to keep your personal information up to date during your move</p>
        </Step>

        <Step
          complete={ordersCompleted}
          completedHeaderText="Orders uploaded"
          editBtnLabel={ordersCompleted ? 'Edit' : ''}
          onEditClick={() => console.log('edit button clicked')}
          headerText="Upload orders"
          actionBtnLabel={!ordersCompleted ? 'Add orders' : ''}
          onActionBtnClick={this.handleCLickToOrders}
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
          onActionBtnClick={this.handleCLickToSelectShipment}
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
          actionBtnDisabled={!hasShipment}
          actionBtnLabel="Review and submit"
          containerClassName="margin-bottom-8"
          headerText="Confirm move request"
          onActionBtnClick={() => console.log('some action')}
          step="4"
        >
          <p className={styles.description}>
            Review your move details and sign the legal paperwork, then send the info on to your move counselor
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
  history: string.isRequired,
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

export default withRouter(connect(mapStateToProps, mapDispatchToProps, mergeProps)(Home));
