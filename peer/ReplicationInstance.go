package peer
import "fmt"
import "errors"
import "github.com/gskartwii/rbxfile"

func DecodeReplicationInstance(isJoinData bool, packet *UDPPacket, context *CommunicationContext, instanceSchema []*InstanceSchemaItem) (*rbxfile.Instance, error) {
	var err error
	thisBitstream := packet.Stream

    referent, err := thisBitstream.ReadObject(isJoinData, context)
	if err != nil {
        return nil, errors.New("while parsing self: " + err.Error())
	}
    thisInstance := context.InstancesByReferent.TryGetInstance(referent)

    schemaIDx, err := thisBitstream.ReadUint16BE()
    if int(schemaIDx) > len(context.StaticSchema.Instances) {
        return thisInstance, errors.New(fmt.Sprintf("class idx %d is higher than %d", schemaIDx, len(context.StaticSchema.Instances)))
    }
    schema := context.StaticSchema.Instances[schemaIDx]
    thisInstance.ClassName = schema.Name
	if DEBUG {
		println("will parse", referent, schema.Name, isJoinData, len(schema.Properties))
	}

    _, err = thisBitstream.ReadBool()
    if err != nil {
        return thisInstance, err
    }
    thisInstance.Properties = make(map[string]rbxfile.Value, len(schema.Properties))

    if isJoinData {
        for i := 0; i < len(schema.Properties); i++ {
            propertyName := schema.Properties[i].Name
			value, err := schema.Properties[i].Decode(ROUND_JOINDATA, packet, context)
			if err != nil {
				return thisInstance, err
			}
			if value != nil {
				thisInstance.Properties[propertyName] = value
			}
        }
    } else {
        for i := 0; i < len(schema.Properties); i++ {
            isStringObject := false
            if  schema.Properties[i].Type == 0x21 ||
                schema.Properties[i].Type == 0x01 ||
                schema.Properties[i].Type == 0x1C ||
                schema.Properties[i].Type == 0x22 ||
                schema.Properties[i].Type == 0x06 ||
                schema.Properties[i].Type == 0x04 ||
                schema.Properties[i].Type == 0x05 ||
                schema.Properties[i].Type == 0x03 {
                    isStringObject = true
            }
            if isStringObject {
                propertyName := schema.Properties[i].Name
				value, err := schema.Properties[i].Decode(ROUND_STRINGS, packet, context)
				if err != nil {
					return thisInstance, err
				}
				if value != nil {
					thisInstance.Properties[propertyName] = value
				}
            }
        }
        for i := 0; i < len(schema.Properties); i++ {
            isStringObject := false
            if  schema.Properties[i].Type == 0x21 ||
                schema.Properties[i].Type == 0x01 ||
                schema.Properties[i].Type == 0x1C ||
                schema.Properties[i].Type == 0x22 ||
                schema.Properties[i].Type == 0x06 ||
                schema.Properties[i].Type == 0x04 ||
                schema.Properties[i].Type == 0x05 ||
                schema.Properties[i].Type == 0x03 {
                    isStringObject = true
            }
            if !isStringObject {
                propertyName := schema.Properties[i].Name
				value, err := schema.Properties[i].Decode(ROUND_OTHER, packet, context)
				if err != nil {
					return thisInstance, err
				}
				if value != nil {
					thisInstance.Properties[propertyName] = value
				}
            }
        }
    }
    referent, err = thisBitstream.ReadObject(isJoinData, context)
    if err != nil {
        return thisInstance, errors.New("while parsing parent: " + err.Error())
    }

    context.InstancesByReferent.AddInstance(Referent(thisInstance.Reference), thisInstance)
    parent := context.InstancesByReferent.TryGetInstance(referent)
    err = parent.AddChild(thisInstance)

    return thisInstance, err
}
