Move Order API Spec

Endpoints
Move  Orders
/move-order       		 						          |          POST        create move order
/move-order/{id} 		 						            |          GET          get a move order
/move-order/{id} 		 						            |          PATCH      update move order
/move-order/{id} 		 						            |          DELETE    delete move order
/move-order/{id}/approve						        |	         PATCH	    approve move order
/move-order/{id}/reject							        |	         PATCH	    reject move order

Modifications for Move  Orders
/move-order/{id}/modification   				    |          POST        create modification for a move order
/move-order/{id}/modification/{id} 			    |          GET          get modification for a move order
/move-order/{id}/modification/{id} 			    |          PATCH      update modification for a move order
/move-order/{id}/modification/{id}  			  |          DELETE    delete modification for a move order
/move-order/{id}/modification/{id}/approve	|	         PATCH	    approve modification for a move order
/move-order/{id}/modification/{id}/reject		|	         PATCH	    reject modification for a move order


MoveOrder
`id`
`payment_method` FK - ENUM table
`move_date` timestamp with time zone
`accounting and appropriation data`?
`origin_duty_station`? FK - will we have a list of these beforehand?
`origin_ppso`?
`destination_duty_station`? FK - will we have a list of these beforehand?
`destination_ppso`?
`type_orders`?
`travel_authorization_name` string
`travel_authorization_date` timestamp with time zone
`travel_authorization_issuing_hq`? FK - will we have a list of these beforehand?
_________________________________
`line_of_accounting`?
`tac`?
`electronic_copy_of_orders`?
`orders_with_all_amendments`?


Modification
`move_order_id` FK
`status` FK - ENUM table  
`management_fee`? is this different per modification? if so, string, or is a payment class which could be FK to enum table
`counseling_fee`? is this different per modification? if so, string, or is a payment class which could be FK to enum table
`transition`?
`hhg_line_items`?
`sit_line_items`?
`sit_pickup`?
`sit_delivery`?
`nts_contractor`?
`nts_pack`?
`ub`?
`shorthaul`?
`local`?
`crating_remarks`?
`shuttle_service_remarks_and_documentation`?
`boat/mobile_home`?
`number`?
`description`?
`quantity`?
`price`?