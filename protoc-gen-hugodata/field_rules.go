package main

import (
	"time"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	"github.com/golang/protobuf/ptypes"
)

type FieldRules struct {
	// Repeated
	MinItems uint64 `yaml:"min_items,omitempty"`
	MaxItems uint64 `yaml:"max_items,omitempty"`
	Unique   bool   `yaml:"unique,omitempty"`
	// Map
	MinPairs uint64 `yaml:"min_pairs,omitempty"`
	MaxPairs uint64 `yaml:"max_pairs,omitempty"`
	NoSparse bool   `yaml:"no_sparse,omitempty"`
	// Enum
	DefinedOnly bool `yaml:"defined_only,omitempty"`
	// Message
	Skip     bool `yaml:"skip,omitempty"`
	Required bool `yaml:"required,omitempty"`
	// String
	Len      uint64 `yaml:"len,omitempty"`
	MinLen   uint64 `yaml:"min_len,omitempty"`
	MaxLen   uint64 `yaml:"max_len,omitempty"`
	LenBytes uint64 `yaml:"len_bytes,omitempty"`
	MinBytes uint64 `yaml:"min_bytes,omitempty"`
	MaxBytes uint64 `yaml:"max_bytes,omitempty"`
	Pattern  string `yaml:"pattern,omitempty"`
	Email    bool   `yaml:"email,omitempty"`
	Hostname bool   `yaml:"hostname,omitempty"`
	URI      bool   `yaml:"uri,omitempty"`
	URIRef   bool   `yaml:"uri_ref,omitempty"`
	Address  bool   `yaml:"address,omitempty"`
	UUID     bool   `yaml:"uuid,omitempty"`
	// String and Bytes
	Prefix      interface{} `yaml:"prefix,omitempty"`
	Suffix      interface{} `yaml:"suffix,omitempty"`
	Contains    interface{} `yaml:"contains,omitempty"`
	NotContains interface{} `yaml:"not_contains,omitempty"`
	IP          bool        `yaml:"ip,omitempty"`
	IPv4        bool        `yaml:"ipv4,omitempty"`
	IPv6        bool        `yaml:"ipv6,omitempty"`
	// Timestamp
	LtNow  bool          `yaml:"lt_now,omitempty"`
	GtNow  bool          `yaml:"gt_now,omitempty"`
	Within time.Duration `yaml:"within,omitempty"`
	// Most types
	Const interface{} `yaml:"const,omitempty"`
	Lt    interface{} `yaml:"lt,omitempty"`
	Lte   interface{} `yaml:"lte,omitempty"`
	Gt    interface{} `yaml:"gt,omitempty"`
	Gte   interface{} `yaml:"gte,omitempty"`
	In    interface{} `yaml:"in,omitempty"`
	NotIn interface{} `yaml:"not_in,omitempty"`
}

func (f *Field) AddFieldRules(src *validate.FieldRules) {
	if src == nil {
		return
	}
	fieldType := f.src.Type()
	if rules := src.GetRepeated(); rules != nil && f.Repeated != nil {
		f.Rules.MinItems = rules.GetMinItems()
		f.Rules.MaxItems = rules.GetMaxItems()
		f.Rules.Unique = rules.GetUnique()
		f.Repeated.Rules.AddFieldRules(rules.GetItems(), fieldType.Element())
	}
	if rules := src.GetMap(); rules != nil && f.MapKey != nil && f.MapValue != nil {
		f.Rules.MinPairs = rules.GetMinPairs()
		f.Rules.MaxPairs = rules.GetMaxPairs()
		f.Rules.NoSparse = rules.GetNoSparse()
		f.MapKey.Rules.AddFieldRules(rules.GetKeys(), fieldType.Key())
		f.MapValue.Rules.AddFieldRules(rules.GetValues(), fieldType.Element())
	}
	f.Rules.AddFieldRules(src, fieldType)
}

func (f *FieldRules) addFloatRules(src *validate.FloatRules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addDoubleRules(src *validate.DoubleRules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addInt32Rules(src *validate.Int32Rules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addInt64Rules(src *validate.Int64Rules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addUint32Rules(src *validate.UInt32Rules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addUint64Rules(src *validate.UInt64Rules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addSint32Rules(src *validate.SInt32Rules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addSint64Rules(src *validate.SInt64Rules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addFixed32Rules(src *validate.Fixed32Rules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addFixed64Rules(src *validate.Fixed64Rules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addSfixed32Rules(src *validate.SFixed32Rules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addSfixed64Rules(src *validate.SFixed64Rules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	if src.Lt != nil {
		f.Lt = src.Lt
	}
	if src.Lte != nil {
		f.Lte = src.Lte
	}
	if src.Gt != nil {
		f.Gt = src.Gt
	}
	if src.Gte != nil {
		f.Gte = src.Gte
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addBoolRules(src *validate.BoolRules) {
	if src.Const != nil {
		f.Const = *src.Const
	}
}

func (f *FieldRules) addStringRules(src *validate.StringRules) {
	if src.Const != nil {
		f.Const = *src.Const
	}
	f.Len = src.GetLen()
	f.MinLen = src.GetMinLen()
	f.MaxLen = src.GetMaxLen()
	f.LenBytes = src.GetLenBytes()
	f.MinBytes = src.GetMinBytes()
	f.MaxBytes = src.GetMaxBytes()
	f.Pattern = src.GetPattern()
	if src.Prefix != nil {
		f.Prefix = src.Prefix
	}
	if src.Suffix != nil {
		f.Suffix = src.Suffix
	}
	if src.Contains != nil {
		f.Contains = src.Contains
	}
	if src.NotContains != nil {
		f.NotContains = src.NotContains
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
	f.IP = src.GetIp()
	f.IPv4 = src.GetIpv4()
	f.IPv6 = src.GetIpv6()
	f.Email = src.GetEmail()
	f.Hostname = src.GetHostname()
	f.URI = src.GetUri()
	f.URIRef = src.GetUriRef()
	f.Address = src.GetAddress()
	f.UUID = src.GetUuid()
}

func (f *FieldRules) addBytesRules(src *validate.BytesRules) {
	if src.Const != nil {
		f.Const = src.Const
	}
	f.Len = src.GetLen()
	f.MinLen = src.GetMinLen()
	f.MaxLen = src.GetMaxLen()
	f.Pattern = src.GetPattern()
	if src.Prefix != nil {
		f.Prefix = src.Prefix
	}
	if src.Suffix != nil {
		f.Suffix = src.Suffix
	}
	if src.Contains != nil {
		f.Contains = src.Contains
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
	f.IP = src.GetIp()
	f.IPv4 = src.GetIpv4()
	f.IPv6 = src.GetIpv6()
}

func (f *FieldRules) addEnumRules(src *validate.EnumRules) {
	f.DefinedOnly = src.GetDefinedOnly()
	if src.Const != nil {
		f.Const = *src.Const
	}
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addMessageRules(src *validate.MessageRules) {
	f.Skip = src.GetSkip()
	f.Required = src.GetRequired()
}

func (f *FieldRules) addAnyRules(src *validate.AnyRules) {
	f.Required = src.GetRequired()
	if src.In != nil {
		f.In = src.In
	}
	if src.NotIn != nil {
		f.NotIn = src.NotIn
	}
}

func (f *FieldRules) addDurationRules(src *validate.DurationRules) {
	if src.Const != nil {
		f.Const, _ = ptypes.Duration(src.Const)
	}
	if src.Lt != nil {
		f.Lt, _ = ptypes.Duration(src.Lt)
	}
	if src.Lte != nil {
		f.Lte, _ = ptypes.Duration(src.Lte)
	}
	if src.Gt != nil {
		f.Gt, _ = ptypes.Duration(src.Gt)
	}
	if src.Gte != nil {
		f.Gte, _ = ptypes.Duration(src.Gte)
	}
	if src.In != nil {
		in := make([]time.Duration, len(src.In))
		for i, p := range src.In {
			in[i], _ = ptypes.Duration(p)
		}
		f.In = in
	}
	if src.NotIn != nil {
		notIn := make([]time.Duration, len(src.NotIn))
		for i, p := range src.NotIn {
			notIn[i], _ = ptypes.Duration(p)
		}
		f.NotIn = notIn
	}
}

func (f *FieldRules) addTimestampRules(src *validate.TimestampRules) {
	f.Required = src.GetRequired()
	if src.Const != nil {
		f.Const, _ = ptypes.Timestamp(src.Const)
	}
	if src.Lt != nil {
		f.Lt, _ = ptypes.Timestamp(src.Lt)
	}
	if src.Lte != nil {
		f.Lte, _ = ptypes.Timestamp(src.Lte)
	}
	if src.Gt != nil {
		f.Gt, _ = ptypes.Timestamp(src.Gt)
	}
	if src.Gte != nil {
		f.Gte, _ = ptypes.Timestamp(src.Gte)
	}
	f.LtNow = src.GetLtNow()
	f.GtNow = src.GetGtNow()
	if src.Within != nil {
		f.Within, _ = ptypes.Duration(src.Within)
	}
}

func (f *FieldRules) AddFieldRules(src *validate.FieldRules, t PGSFieldType) {
	if src == nil {
		return
	}
	if rules := src.GetFloat(); rules != nil {
		f.addFloatRules(rules)
	}
	if rules := src.GetDouble(); rules != nil {
		f.addDoubleRules(rules)
	}
	if rules := src.GetInt32(); rules != nil {
		f.addInt32Rules(rules)
	}
	if rules := src.GetInt64(); rules != nil {
		f.addInt64Rules(rules)
	}
	if rules := src.GetUint32(); rules != nil {
		f.addUint32Rules(rules)
	}
	if rules := src.GetUint64(); rules != nil {
		f.addUint64Rules(rules)
	}
	if rules := src.GetSint32(); rules != nil {
		f.addSint32Rules(rules)
	}
	if rules := src.GetSint64(); rules != nil {
		f.addSint64Rules(rules)
	}
	if rules := src.GetFixed32(); rules != nil {
		f.addFixed32Rules(rules)
	}
	if rules := src.GetFixed64(); rules != nil {
		f.addFixed64Rules(rules)
	}
	if rules := src.GetSfixed32(); rules != nil {
		f.addSfixed32Rules(rules)
	}
	if rules := src.GetSfixed64(); rules != nil {
		f.addSfixed64Rules(rules)
	}
	if rules := src.GetBool(); rules != nil {
		f.addBoolRules(rules)
	}
	if rules := src.GetString_(); rules != nil {
		f.addStringRules(rules)
	}
	if rules := src.GetBytes(); rules != nil {
		f.addBytesRules(rules)
	}
	if rules := src.GetEnum(); rules != nil && t.IsEnum() {
		f.addEnumRules(rules)
	}
	if rules := src.GetMessage(); rules != nil && t.IsEmbed() {
		f.addMessageRules(rules)
	}
	if rules := src.GetAny(); rules != nil {
		f.addAnyRules(rules)
	}
	if rules := src.GetDuration(); rules != nil {
		f.addDurationRules(rules)
	}
	if rules := src.GetTimestamp(); rules != nil {
		f.addTimestampRules(rules)
	}
}
