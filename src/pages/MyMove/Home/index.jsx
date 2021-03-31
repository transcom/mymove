import React, { Component } from 'react';
import { arrayOf, bool, shape, string, node, func } from 'prop-types';
import moment from 'moment';
import { connect } from 'react-redux';
import { Button } from '@trussworks/react-uswds';
import { generatePath } from 'react-router';

import styles from './Home.module.scss';
import {
  HelperNeedsOrders,
  HelperNeedsShipment,
  HelperNeedsSubmitMove,
  HelperSubmittedMove,
  HelperSubmittedPPM,
} from './HomeHelpers';

import ScrollToTop from 'components/ScrollToTop';
import { customerRoutes } from 'constants/routes';
import { withContext } from 'shared/AppContext';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Step from 'components/Customer/Home/Step';
import DocsUploaded from 'components/Customer/Home/DocsUploaded';
import ShipmentList from 'components/Customer/Home/ShipmentList';
import Contact from 'components/Customer/Home/Contact';
import SectionWrapper from 'components/Customer/SectionWrapper';
import PrintableLegalese from 'components/Customer/Home/PrintableLegalese';
import {
  selectServiceMemberFromLoggedInUser,
  selectIsProfileComplete,
  selectCurrentOrders,
  selectCurrentMove,
  selectCurrentPPM,
  selectMTOShipmentsForCurrentMove,
  selectUploadsForCurrentOrders,
} from 'store/entities/selectors';
import {
  getSignedCertification as getSignedCertificationAction,
  selectSignedCertification,
} from 'shared/Entities/modules/signed_certifications';
import { SHIPMENT_OPTIONS, MOVE_STATUSES } from 'shared/constants';
import { formatCustomerDate } from 'utils/formatters';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import { MtoShipmentShape, UploadShape, HistoryShape, MoveShape, OrdersShape } from 'types/customerShapes';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';

const Description = ({ className, children, dataTestId }) => (
  <p className={`${styles.description} ${className}`} data-testid={dataTestId}>
    {children}
  </p>
);

Description.propTypes = {
  className: string,
  children: node.isRequired,
  dataTestId: string,
};

Description.defaultProps = {
  className: '',
  dataTestId: '',
};

export class Home extends Component {
  componentDidMount() {
    const { move, getSignedCertification } = this.props;
    if (Object.entries(move).length && move.status === MOVE_STATUSES.SUBMITTED) {
      getSignedCertification(move.id);
    }
  }

  componentDidUpdate(prevProps) {
    const { move, getSignedCertification } = this.props;

    if (!Object.entries(prevProps.move).length && Object.entries(move).length) {
      getSignedCertification(move.id);
    }
  }

  get hasOrders() {
    const { orders, uploadedOrderDocuments } = this.props;
    return !!Object.keys(orders).length && !!uploadedOrderDocuments.length;
  }

  get hasOrdersNoUpload() {
    const { orders, uploadedOrderDocuments } = this.props;
    return !!Object.keys(orders).length && !uploadedOrderDocuments.length;
  }

  get hasAnyShipments() {
    const { mtoShipments, currentPpm } = this.props;
    return (this.hasOrders && !!mtoShipments.length) || !!Object.keys(currentPpm).length;
  }

  get hasSubmittedMove() {
    const { move } = this.props;
    return !!Object.keys(move).length && move.status !== 'DRAFT';
  }

  get hasHHGShipment() {
    const { mtoShipments } = this.props;
    return mtoShipments.some((s) => s.shipmentType === SHIPMENT_OPTIONS.HHG);
  }

  get hasNTSShipment() {
    const { mtoShipments } = this.props;
    return mtoShipments.some((s) => s.shipmentType === SHIPMENT_OPTIONS.NTS);
  }

  get hasPPMShipment() {
    const { currentPpm } = this.props;
    return !!Object.keys(currentPpm).length;
  }

  get shipmentActionBtnLabel() {
    if (this.hasSubmittedMove && this.hasPPMShipment) {
      return '';
    }
    if (this.hasAnyShipments) {
      return 'Add another shipment';
    }
    return 'Plan your shipments';
  }

  renderHelper = () => {
    if (!this.hasOrders) return <HelperNeedsOrders />;
    if (!this.hasAnyShipments) return <HelperNeedsShipment />;
    if (!this.hasSubmittedMove) return <HelperNeedsSubmitMove />;
    if (this.hasPPMShipment)
      return (
        <>
          <HelperSubmittedMove />
          <HelperSubmittedPPM />
        </>
      );
    return <HelperSubmittedMove />;
  };

  renderCustomerHeaderText = () => {
    const { serviceMember, orders, move } = this.props;
    return (
      <>
        <p>
          You’re moving to <strong>{orders.new_duty_station.name}</strong> from{' '}
          <strong>{orders.origin_duty_station?.name}.</strong> Report by{' '}
          <strong>{moment(orders.report_by_date).format('DD MMM YYYY')}.</strong>
        </p>

        <dl className={styles.subheaderContainer}>
          <div className={styles.subheaderSubsection}>
            <dt>Weight allowance</dt>
            <dd>{serviceMember.weight_allotment.total_weight_self} lbs.</dd>
          </div>
          {move.locator && (
            <div className={styles.subheaderSubsection}>
              <dt>Move code</dt>
              <dd>#{move.locator}</dd>
            </div>
          )}
        </dl>
      </>
    );
  };

  handleShipmentClick = (shipmentId, shipmentNumber, shipmentType) => {
    const { move, history } = this.props;
    let queryString = '';
    if (shipmentNumber) {
      queryString = `?shipmentNumber=${shipmentNumber}`;
    }

    let destLink = '';
    if (shipmentType === 'PPM') {
      destLink = `/moves/${move.id}/review/edit-date-and-location`;
    } else if (shipmentType === 'HHG') {
      destLink = `${generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: move.id,
        mtoShipmentId: shipmentId,
      })}${queryString}`;
    } else {
      // nts/ntsr shipment
      destLink = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: move.id,
        mtoShipmentId: shipmentId,
      });
    }

    history.push(destLink);
  };

  handleNewPathClick = (path) => {
    const { history } = this.props;
    history.push(path);
  };

  sortAllShipments = (mtoShipments, currentPpm) => {
    const allShipments = JSON.parse(JSON.stringify(mtoShipments));
    if (Object.keys(currentPpm).length) {
      const ppm = JSON.parse(JSON.stringify(currentPpm));
      ppm.shipmentType = SHIPMENT_OPTIONS.PPM;
      // workaround for differing cases between mtoShipments and ppms (bigger change needed on yaml)
      ppm.createdAt = ppm.created_at;
      delete ppm.created_at;

      allShipments.push(ppm);
    }
    allShipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));

    return allShipments;
  };

  handlePrintLegalese = (e) => {
    e.preventDefault();
    window.print();
  };

  render() {
    const {
      currentPpm,
      isProfileComplete,
      move,
      mtoShipments,
      serviceMember,
      signedCertification,
      uploadedOrderDocuments,
    } = this.props;

    // early return if loading user/service member
    if (!serviceMember) {
      return (
        <div className={styles.homeContainer}>
          <div className={`usa-prose grid-container ${styles['grid-container']}`}>
            <LoadingPlaceholder />
          </div>
        </div>
      );
    }

    // eslint-disable-next-line camelcase
    const { current_station } = serviceMember;
    const ordersPath = this.hasOrdersNoUpload ? customerRoutes.ORDERS_UPLOAD_PATH : customerRoutes.ORDERS_INFO_PATH;

    const shipmentSelectionPath =
      move?.id &&
      (this.hasAnyShipments
        ? generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, { moveId: move.id })
        : generatePath(customerRoutes.SHIPMENT_MOVING_INFO_PATH, { moveId: move.id }));

    const confirmationPath = move?.id && generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId: move.id });
    const profileEditPath = '/moves/review/edit-profile';
    const ordersEditPath = `/moves/${move.id}/review/edit-orders`;
    const allSortedShipments = this.sortAllShipments(mtoShipments, currentPpm);

    return (
      <>
        <ScrollToTop />
        <div className={styles.homeContainer}>
          <header data-testid="customer-header" className={styles['customer-header']}>
            <div className={`usa-prose grid-container ${styles['grid-container']}`}>
              <h2>
                {serviceMember.first_name} {serviceMember.last_name}
              </h2>
              {(this.hasOrdersNoUpload || this.hasOrders) && this.renderCustomerHeaderText()}
            </div>
          </header>
          <div className={`usa-prose grid-container ${styles['grid-container']}`}>
            <ConnectedFlashMessage />

            {isProfileComplete && (
              <>
                {this.renderHelper()}
                <SectionWrapper>
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
                    editBtnLabel={this.hasOrders && !this.hasSubmittedMove ? 'Edit' : ''}
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
                    actionBtnLabel={this.shipmentActionBtnLabel}
                    actionBtnDisabled={!this.hasOrders || (this.hasSubmittedMove && this.doesPpmAlreadyExist)}
                    actionBtnId="shipment-selection-btn"
                    onActionBtnClick={() => this.handleNewPathClick(shipmentSelectionPath)}
                    complete={this.hasAnyShipments}
                    completedHeaderText="Shipments"
                    headerText="Shipment selection"
                    secondaryBtn={this.hasAnyShipments}
                    secondaryClassName="margin-top-2"
                    step="3"
                  >
                    {this.hasAnyShipments ? (
                      <div>
                        {this.hasSubmittedMove && !this.doesPpmAlreadyExist && (
                          <p className={styles.descriptionExtra}>If you need to add shipments, let your movers know.</p>
                        )}
                        <ShipmentList
                          shipments={allSortedShipments}
                          onShipmentClick={this.handleShipmentClick}
                          moveSubmitted={this.hasSubmittedMove}
                        />
                      </div>
                    ) : (
                      <Description>
                        Tell us where you&apos;re going and when you want to get there. We&apos;ll help you set up
                        shipments to make it work.
                      </Description>
                    )}
                  </Step>
                  <Step
                    actionBtnDisabled={!this.hasAnyShipments}
                    actionBtnId="review-and-submit-btn"
                    actionBtnLabel={!this.hasSubmittedMove ? 'Review and submit' : 'Review your request'}
                    complete={this.hasSubmittedMove}
                    completedHeaderText="Move request confirmed"
                    containerClassName="margin-bottom-8"
                    headerText="Confirm move request"
                    onActionBtnClick={() => this.handleNewPathClick(confirmationPath)}
                    secondaryBtn={this.hasSubmittedMove}
                    secondaryBtnClassName={styles.secondaryBtn}
                    step="4"
                  >
                    {this.hasSubmittedMove ? (
                      <Description className={styles.moveSubmittedDescription} dataTestId="move-submitted-description">
                        Move submitted {formatCustomerDate(move.submitted_at)}.<br />
                        <Button unstyled onClick={this.handlePrintLegalese} className={styles.printBtn}>
                          Print the legal agreement
                        </Button>
                      </Description>
                    ) : (
                      <Description>
                        Review your move details and sign the legal paperwork, then send the info on to your move
                        counselor.
                      </Description>
                    )}
                  </Step>
                </SectionWrapper>
                <Contact
                  header="Contacts"
                  dutyStationName={current_station?.transportation_office?.name}
                  officeType="Origin Transportation Office"
                  telephone={current_station?.transportation_office?.phone_lines[0]}
                  moveSubmitted={move.status === MOVE_STATUSES.SUBMITTED}
                />
              </>
            )}
          </div>
        </div>
        <PrintableLegalese signature={signedCertification.signature} signatureDate={signedCertification.created_at} />
      </>
    );
  }
}

Home.propTypes = {
  orders: OrdersShape,
  serviceMember: shape({
    first_name: string,
    last_name: string,
  }),
  mtoShipments: arrayOf(MtoShipmentShape).isRequired,
  currentPpm: shape({
    id: string,
    shipmentType: string,
  }).isRequired,
  uploadedOrderDocuments: arrayOf(UploadShape).isRequired,
  history: HistoryShape.isRequired,
  move: MoveShape.isRequired,
  isProfileComplete: bool.isRequired,
  signedCertification: shape({
    signature: string,
    created_at: string,
  }),
  getSignedCertification: func.isRequired,
};

Home.defaultProps = {
  orders: null,
  serviceMember: null,
  signedCertification: {},
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const move = selectCurrentMove(state) || {};

  return {
    currentPpm: selectCurrentPPM(state) || {},
    isProfileComplete: selectIsProfileComplete(state),
    orders: selectCurrentOrders(state) || {},
    uploadedOrderDocuments: selectUploadsForCurrentOrders(state),
    serviceMember,
    backupContacts: serviceMember?.backup_contacts || [],
    signedCertification: selectSignedCertification(state),
    // TODO: change when we support PPM shipments as well
    mtoShipments: selectMTOShipmentsForCurrentMove(state),
    // TODO: change when we support multiple moves
    move,
  };
};

const mapDispatchToProps = {
  getSignedCertification: getSignedCertificationAction,
};

// in order to avoid setting up proxy server only for storybook, pass in stub function so API requests don't fail
const mergeProps = (stateProps, dispatchProps, ownProps) => ({
  ...stateProps,
  ...dispatchProps,
  ...ownProps,
});

export default withContext(
  connect(
    mapStateToProps,
    mapDispatchToProps,
    mergeProps,
  )(requireCustomerState(Home, profileStates.BACKUP_CONTACTS_COMPLETE)),
);
