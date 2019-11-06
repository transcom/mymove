import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import ppmCar from 'scenes/Landing/images/ppm-car.svg';
import PPMStatusTimeline from 'scenes/Landing/PPMStatusTimeline';
import PpmMoveDetails from 'scenes/Landing/MoveSummary/PpmMoveDetails';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import { get } from 'lodash';

const ApprovedMoveSummary = ({ ppm, move, weightTicketSets, isMissingWeightTicketDocuments, incentiveEstimate }) => {
  const paymentRequested = ppm.status === 'PAYMENT_REQUESTED';
  const ppmPaymentRequestIntroRoute = `moves/${move.id}/ppm-payment-request-intro`;
  const ppmPaymentRequestReviewRoute = `moves/${move.id}/ppm-payment-review`;
  return (
    <div>
      <div className="shipment_box">
        <div className="shipment_type">
          <img className="move_sm" src={ppmCar} alt="ppm-car" />
          Move your own stuff (PPM)
        </div>

        <div className="shipment_box_contents">
          <PPMStatusTimeline ppm={ppm} />
          <div className="step-contents">
            <div className="status_box usa-width-two-thirds">
              {paymentRequested ? (
                isMissingWeightTicketDocuments ? (
                  <div className="step">
                    <div className="title">Next step: Contact the PPPO office</div>
                    <div>
                      You will need to go into the PPPO office in order to take care of your missing weight ticket.
                    </div>
                    <Link
                      data-cy="edit-payment-request"
                      to={ppmPaymentRequestReviewRoute}
                      className="usa-button usa-button-secondary"
                    >
                      Edit Payment Request
                    </Link>
                  </div>
                ) : (
                  <div className="step">
                    <div className="title">What's next?</div>
                    <div>
                      We'll email you a link so you can see and download your final payment paperwork.
                      <br />
                      <br />
                      We've also sent your paperwork to Finance. They'll review it, determine a final amount, then send
                      your payment.
                    </div>
                    <Link
                      data-cy="edit-payment-request"
                      to={ppmPaymentRequestReviewRoute}
                      className="usa-button usa-button-secondary"
                    >
                      Edit Payment Request
                    </Link>
                  </div>
                )
              ) : (
                <div className="step">
                  {weightTicketSets.length ? (
                    <>
                      <div className="title">Next Step: Finish requesting payment</div>
                      <div>Continue uploading your weight tickets and expense to get paid after your move is done.</div>
                      <Link to={ppmPaymentRequestReviewRoute} className="usa-button usa-button-secondary">
                        Continue Requesting Payment
                      </Link>
                    </>
                  ) : (
                    <>
                      <div className="title">Next Step: Request payment</div>
                      <div>
                        Request a PPM payment, a storage payment, or an advance against your PPM payment before your
                        move is done.
                      </div>
                      <Link to={ppmPaymentRequestIntroRoute} className="usa-button usa-button-secondary">
                        Request Payment
                      </Link>
                    </>
                  )}
                </div>
              )}
            </div>
            <div className="usa-width-one-third">
              <PpmMoveDetails ppm={ppm} />
            </div>
          </div>
          <div className="step-links" />
        </div>
      </div>
    </div>
  );
};

const mapStateToProps = (state, { move }) => ({
  weightTicketSets: selectPPMCloseoutDocumentsForMove(state, move.id, ['WEIGHT_TICKET_SET']),
  incentiveEstimate: get(state, 'ppm.incentive_estimate_min'),
});

export default connect(mapStateToProps)(ApprovedMoveSummary);
