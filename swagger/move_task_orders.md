Move Task Order Spec


Endpoints
Move Task Orders
/move-task-order       		 						                  |          POST        create move task order
/move-task-order/{id} 		 						                  |          GET          get a move task order
/move-task-order/{id} 		 						                  |          PATCH      update move task order
/move-task-order/{id} 		 						                  |          DELETE    delete move task order
/move-task-order/{id}/approve						                |	         PATCH	    approve move task order
/move-task-order/{id}/reject							              |	         PATCH	    reject move task order

Shipment Line Items for Move Task Orders
/move-task-order/{id}/shipment-line-items   				    |          POST        create shipment line item for a move task order
/move-task-order/{id}/shipment-line-items/{id} 			    |          GET          get shipment line item for a move task order
/move-task-order/{id}/shipment-line-items/{id} 			    |          PATCH      update shipment line item for a move task order
/move-task-order/{id}/shipment-line-items/{id}  			  |          DELETE    delete shipment line item for a move task order
/move-task-order/{id}/shipment-line-items/{id}/approve	|	         PATCH	    approve shipment line item for a move task order
/move-task-order/{id}/shipment-line-items/{id}/reject		|	         PATCH	    reject shipment line item for a move task order

MoveTaskOrder
`id`
`move_date`


ShipmentLineItem
`move_task_order_id`
`status`