import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Alert, CardGroup, CardHeader, CardBody, CardMedia, GridContainer, Grid } from '@trussworks/react-uswds';
import { func, number, string } from 'prop-types';
import { generatePath } from 'react-router';

import { generalRoutes, customerRoutes } from 'constants/routes';
import ScrollToTop from 'components/ScrollToTop';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import PPMShipmentImg from 'images/ppm-shipment.png';
import HHGShipmentImg from 'images/hhg-shipment.jpg';
import MoveCounselorImg from 'images/move-counselor.jpg';
import MovingTruckImg from 'images/moving-truck.jpg';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MilmoveCard from 'components/Customer/MilmoveCard/MilmoveCard';
import { fetchLatestOrders as fetchLatestOrdersAction } from 'shared/Entities/modules/orders';
import { formatWeight } from 'shared/formatters';
import { selectServiceMemberFromLoggedInUser, selectCurrentOrders } from 'store/entities/selectors';
import { RouteProps } from 'types/router';

export class MovingInfo extends Component {
  componentDidMount() {
    const { serviceMemberId, fetchLatestOrders } = this.props;
    fetchLatestOrders(serviceMemberId);
  }

  render() {
    const { entitlementWeight, history, match } = this.props;

    const nextPath = generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, {
      moveId: match.params.moveId,
    });

    return (
      <GridContainer>
        <ScrollToTop />
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <h1 data-testid="shipmentsHeader">Tips for planning your shipments</h1>
            {entitlementWeight === 0 ? (
              <></>
            ) : (
              <Alert
                data-testid="shipmentsAlert"
                type="info"
                heading={`You can move ${formatWeight(entitlementWeight)} for free`}
                noIcon
              >
                The government will pay to move that much weight. Your whole move, no matter how many shipments it
                takes.
                <br />
                <br />
                If you move more weight, you’ll need to pay for the excess. We’ll tell you if it looks like that could
                happen.
              </Alert>
            )}
            <SectionWrapper>
              <CardGroup>
                <MilmoveCard>
                  <CardHeader>
                    <h3 data-testid="shipmentsSubHeader">Hold on to things you’ll need quickly</h3>
                  </CardHeader>
                  <CardMedia inset>
                    <img src={PPMShipmentImg} alt="PPM Shipment" />
                  </CardMedia>
                  <CardBody>
                    <p>Hand-carry important documents — ID, medical info, orders, school records, etc.</p>
                    <p>
                      Pack a set of things that you’ll need when you arrive — clothes, electronics, chargers, cleaning
                      supplies, etc. Valuables that can’t be replaced are also a good idea.
                    </p>
                    <p>To be paid for moving these things, select a PPM shipment.</p>
                  </CardBody>
                </MilmoveCard>
                <MilmoveCard>
                  <CardHeader>
                    <h3 data-testid="shipmentsSubHeader">One move, several parts</h3>
                  </CardHeader>
                  <CardMedia inset>
                    <img src={HHGShipmentImg} alt="HHG Shipment" />
                  </CardMedia>
                  <CardBody>
                    <p>
                      It’s common to move a few things yourself and have professional movers pack and move the rest.
                    </p>
                    <p>
                      You can have things picked up or delivered to more than one place — your home and an office, for
                      example. But multiple shipments make it easier to go over weight and end up paying for part of
                      your move yourself.
                    </p>
                  </CardBody>
                </MilmoveCard>
                <MilmoveCard>
                  <CardHeader>
                    <h3 data-testid="shipmentsSubHeader">Talk to your move counselor</h3>
                  </CardHeader>
                  <CardMedia inset>
                    <img src={MoveCounselorImg} alt="Move counselor" />
                  </CardMedia>
                  <CardBody>
                    <p>
                      A session with a move counselor is free. Counselors have a lot of experience with military moves
                      and can steer you through complicated situations.
                    </p>
                    <p>Your counselor can identify:</p>
                    <ul>
                      <li>belongings that won’t count against your weight allowance</li>
                      <li>excess weight, excess distance, and other things that can cost you money</li>
                      <li>things to make your move easier</li>
                    </ul>
                  </CardBody>
                </MilmoveCard>
                <MilmoveCard>
                  <CardHeader>
                    <h3 data-testid="shipmentsSubHeader">Talk to your movers</h3>
                  </CardHeader>
                  <CardMedia inset>
                    <img src={MovingTruckImg} alt="Moving truck" />
                  </CardMedia>
                  <CardBody>
                    <p>
                      If you have any shipments using professional movers, you’ll be referred to a point of contact for
                      your move.
                    </p>
                    <p>When things get complicated or you have questions during your move, they are there to help.</p>
                    <p>
                      It’s OK if things change after you submit your move info. Your movers or your counselor will make
                      things work.
                    </p>
                  </CardBody>
                </MilmoveCard>
              </CardGroup>
            </SectionWrapper>

            <WizardNavigation
              isFirstPage
              showFinishLater
              onNextClick={() => {
                history.push(nextPath);
              }}
              onCancelClick={() => {
                history.push(generalRoutes.HOME_PATH);
              }}
            />
          </Grid>
        </Grid>
      </GridContainer>
    );
  }
}

MovingInfo.propTypes = {
  entitlementWeight: number,
  fetchLatestOrders: func.isRequired,
  serviceMemberId: string.isRequired,
  ...RouteProps,
};

MovingInfo.defaultProps = {
  entitlementWeight: 0,
};

function mapStateToProps(state) {
  const orders = selectCurrentOrders(state);
  const entitlementWeight = orders.authorizedWeight;
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;

  return {
    entitlementWeight,
    serviceMemberId,
  };
}

const mapDispatchToProps = {
  fetchLatestOrders: fetchLatestOrdersAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(MovingInfo);
