/* eslint-disable react/no-unused-prop-types */
/* eslint-disable no-console */
import React, { Component } from 'react';
import { func, arrayOf, shape, string, node } from 'prop-types';
import moment from 'moment';
import { connect } from 'react-redux';

import styles from './Home.module.scss';

import Helper from 'components/Customer/Home/Helper';
import Step from 'components/Customer/Home/Step';
import DocsUploaded from 'components/Customer/Home/DocsUploaded';
import ShipmentList from 'components/Customer/Home/ShipmentList';
import Contact from 'components/Customer/Home/Contact';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { selectServiceMemberFromLoggedInUser } from 'shared/Entities/modules/serviceMembers';
import { selectUploadedOrders, selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { selectActiveOrLatestMove } from 'shared/Entities/modules/moves';
import {
  selectMTOShipmentsByMoveId,
  loadMTOShipments as loadMTOShipmentsAction,
} from 'shared/Entities/modules/mtoShipments';

const Description = ({ children }) => <p className={styles.description}>{children}</p>;

Description.propTypes = {
  children: node.isRequired,
};

class Home extends Component {
  componentDidMount() {
    const { showLoggedInUser, move, loadMTOShipments } = this.props;
    showLoggedInUser();
    if (move.id) {
      loadMTOShipments(move.id);
    }
  }

  get hasOrders() {
    const { orders, uploadedOrderDocuments } = this.props;
    return !!Object.keys(orders).length && !!uploadedOrderDocuments.length;
  }

  get hasShipment() {
    const { shipments } = this.props;
    // TODO: check for PPM when PPM is integrated
    return this.hasOrders && !!shipments.length;
  }

  get hasSubmittedMove() {
    const { move } = this.props;
    return !!Object.keys(move).length && move.status !== 'DRAFT';
  }

  get getHelperHeaderText() {
    if (!this.hasOrders) {
      return 'Next step: Add your orders';
    }

    if (!this.hasShipment) {
      return 'Gather this info, then plan your shipments';
    }

    if (this.hasShipment && !this.hasSubmittedMove) {
      return 'Time to submit your move';
    }

    if (this.hasSubmittedMove) {
      return 'Track your HHG move here';
    }

    return '';
  }

  renderHelperListItems = (helperList) => {
    return helperList.map((listItemText) => (
      <li key={listItemText}>
        <span>{listItemText}</span>
      </li>
    ));
  };

  renderHelperDescription = () => {
    if (!this.hasOrders) {
      return (
        <ul>
          {this.renderHelperListItems([
            'If you have a hard copy, you can take photos of each page',
            'If you have a PDF, you can upload that',
          ])}
        </ul>
      );
    }

    if (!this.hasShipment) {
      return (
        <ul>
          {this.renderHelperListItems([
            'Preferred moving details',
            'Destination address (your new place, your duty station ZIP, or somewhere else)',
            'Names and contact info for anyone you authorize to act on your behalf',
          ])}
        </ul>
      );
    }

    if (this.hasShipment && !this.hasSubmittedMove) {
      return (
        <ul>
          {this.renderHelperListItems([
            "Double check the info you've entered",
            'Sign the legal agreement',
            "You'll hear from a move counselor or your transportation office within a few days",
          ])}
        </ul>
      );
    }

    if (this.hasSubmittedMove) {
      return (
        <ul>
          {this.renderHelperListItems([
            'Create a custom checklist at Plan My Move',
            'Learn more about your new duty station',
          ])}
        </ul>
      );
    }

    return null;
  };

  renderCustomerHeader = () => {
    const { serviceMember, orders } = this.props;
    if (!this.hasOrders) {
      return (
        <p>
          You&apos;re leaving <strong>{serviceMember?.['current_station']?.name}</strong>
        </p>
      );
    }
    return (
      <p>
        You&apos;re moving to <strong>{orders.new_duty_station.name}</strong> from{' '}
        <strong>{serviceMember.current_station.name}.</strong> Report by{' '}
        <strong>{moment(orders.report_by_date).format('DD MMM YYYY')}.</strong>
        <br />
        Weight allowance: <strong>{serviceMember.weight_allotment.total_weight_self} lbs</strong>
      </p>
    );
  };

  handleShipmentClick = (shipment) => {
    // TODO: use shipment id in review path with multiple shipments functionality
    console.log('this is the shipment', shipment);
    const { history } = this.props;
    history.push('/moves/review/edit-shipment');
  };

  handleNewPathClick = (path) => {
    const { history } = this.props;
    history.push(path);
  };

  render() {
    const { move, serviceMember, uploadedOrderDocuments, shipments } = this.props;
    const ordersPath = '/orders/';
    const shipmentSelectionPath = `/moves/${move.id}/select-type`;
    const confirmationPath = `/moves/${move.id}/review`;
    const profileEditPath = '/moves/review/edit-profile';
    const ordersEditPath = `/moves/${move.id}/review/edit-orders`;

    return (
      <div className={`usa-prose grid-container ${styles['grid-container']}`}>
        <header data-testid="customer-header" className={styles['customer-header']}>
          <h2>
            {serviceMember?.['first_name']} {serviceMember?.['last_name']}
          </h2>
          {this.renderCustomerHeader()}
        </header>
        <Helper title={this.getHelperHeaderText}>{this.renderHelperDescription()}</Helper>
        <Step
          complete={serviceMember.is_profile_complete}
          completedHeaderText="Profile complete"
          editBtnLabel="Edit"
          headerText="Profile complete"
          step="1"
          onEditBtnClick={() => this.handleNewPathClick(profileEditPath)}
        >
          <Description>Make sure to keep your personal information up to date during your move</Description>
        </Step>

        <Step
          complete={this.hasOrders}
          completedHeaderText="Orders uploaded"
          editBtnLabel={this.hasOrders ? 'Edit' : ''}
          onEditBtnClick={() => this.handleNewPathClick(ordersEditPath)}
          headerText="Upload orders"
          actionBtnLabel={!this.hasOrders ? 'Add orders' : ''}
          onActionBtnClick={() => this.handleNewPathClick(ordersPath)}
          step="2"
        >
          {this.hasOrders ? (
            <DocsUploaded files={uploadedOrderDocuments} />
          ) : (
            <Description>Upload photos of each page, or upload a PDF.</Description>
          )}
        </Step>

        <Step
          actionBtnLabel={this.hasShipment ? 'Add another shipment' : 'Plan your shipments'}
          actionBtnDisabled={!this.hasOrders}
          onActionBtnClick={() => this.handleNewPathClick(shipmentSelectionPath)}
          complete={this.hasShipment}
          completedHeaderText="Shipments"
          headerText="Shipment selection"
          secondaryBtn={this.hasShipment}
          secondaryClassName="margin-top-2"
          step="3"
        >
          {this.hasShipment ? (
            <ShipmentList shipments={shipments} onShipmentClick={this.handleShipmentClick} />
          ) : (
            <Description>
              Tell us where you&apos;re going and when you want to get there. We&apos;ll help you set up shipments to
              make it work.
            </Description>
          )}
        </Step>

        <Step
          complete={this.hasSubmittedMove}
          actionBtnDisabled={!this.hasShipment}
          actionBtnLabel={!this.hasSubmittedMove ? 'Review and submit' : ''}
          containerClassName="margin-bottom-8"
          headerText="Confirm move request"
          completedHeaderText="Move request confirmed"
          onActionBtnClick={() => this.handleNewPathClick(confirmationPath)}
          step="4"
        >
          <p className={styles.description}>
            {this.hasSubmittedMove
              ? 'Move submitted.'
              : 'Review your move details and sign the legal paperwork, then send the info on to your move counselor.'}
          </p>
          {this.hasSubmittedMove && (
            <Description>
              Review your move details and sign the legal paperwork, then send the info on to your move counselor
            </Description>
          )}
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
  orders: shape({}).isRequired,
  serviceMember: shape({
    first_name: string.isRequired,
    last_name: string.isRequired,
  }).isRequired,
  showLoggedInUser: func.isRequired,
  loadMTOShipments: func.isRequired,
  shipments: arrayOf(
    shape({
      id: string,
      shipmentType: string,
    }),
  ).isRequired,
  uploadedOrderDocuments: arrayOf(
    shape({
      filename: string.isRequired,
    }),
  ).isRequired,
  history: shape({}).isRequired,
  move: shape({}).isRequired,
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const move = selectActiveOrLatestMove(state);
  return {
    orders: selectActiveOrLatestOrdersFromEntities(state),
    uploadedOrderDocuments: selectUploadedOrders(state),
    serviceMember,
    // TODO: change when we support PPM shipments as well
    shipments: selectMTOShipmentsByMoveId(state, move.id),
    // TODO: change when we support multiple moves
    move,
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
  loadMTOShipments: loadMTOShipmentsAction,
};

export default connect(mapStateToProps, mapDispatchToProps, mergeProps)(Home);
